package mail

import (
	"context"
	"errors"
	"fmt"
	"net/smtp"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"

	mail "github.com/jordan-wright/email"
)

type Mailer struct {
	From string

	pool *mail.Pool
	maxEmailPerDay int32
	timeout time.Duration

	errChan chan error
	doneChan chan struct{}
}

type Message struct {
	To string
	Subject string
	PlainContent string
	HTMLContent string
}

func NewMailer(host, address, email, password string) *Mailer {
	p, err := mail.NewPool(host, 4, smtp.PlainAuth("", email, password, address))
	if err != nil {
		log.Fatal().Msg("Error on creating mailer: " + err.Error())
	}

	mailer := &Mailer{
		From: email,
		pool: p, 
		maxEmailPerDay: 0, 
		timeout: time.Second * 10,
		doneChan: make(chan struct{}), 
		errChan: make(chan error),
	}

	log.Info().Msg(fmt.Sprintf("Creating new mail on host : %s | email : %s | password : %s", host, email, password))
	return mailer
}

func (m *Mailer) Send(ctx context.Context, msg Message) {
	e := mail.NewEmail()

	e.From = m.From
	e.To = []string{msg.To}
	e.Subject = msg.Subject
	e.HTML = []byte(msg.HTMLContent)

	go m.send(ctx, e)
}

func (m *Mailer) send(ctx context.Context, msg *mail.Email) {
	select {
	case <-ctx.Done():
		return

	default:
		if m.maxEmailPerDay == 1500 {
			m.errChan <- errors.New("cant send email right now, limit reached")
			return
		}

		if msg == nil {
			m.errChan <- errors.New("email that want to send is nil")
			return
		}

		err := m.pool.Send(msg, m.timeout)
		if err != nil {
			m.errChan <- err
			return
		}

		atomic.AddInt32(&m.maxEmailPerDay, 1)
		return
	}
}


func (m *Mailer) Listen() {
	for {
		select{
		case <- m.doneChan:
			close(m.errChan)
			close(m.doneChan)
			return

		case err := <-m.errChan:
			if err != nil {
				log.Error().Msg("Error on mailer: " + err.Error())
			}
		}
	}
}

func (m *Mailer) Close() { 
	m.doneChan <- struct{}{} 
	m.pool.Close() 
}