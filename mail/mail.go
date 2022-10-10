package mail

import (
	"context"
	"net/smtp"
	"time"

	"github.com/rs/zerolog/log"

	mail "github.com/jordan-wright/email"
)

type Mailer struct {
	From string

	pool *mail.Pool
	emailQueue chan *mail.Email
	maxEmailPerDay int16
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
		log.Fatal().Msg("Error on closing database" + err.Error())
	}

	mailer := &Mailer{
		From: email,
		pool: p, 
		maxEmailPerDay: 1500, 
		emailQueue: make(chan *mail.Email),
		timeout: time.Second * 10,
		doneChan: make(chan struct{}), 
		errChan: make(chan error),
	}
	return mailer
}

func (m *Mailer) Send(ctx context.Context, msg Message) {
	e := mail.NewEmail()

	e.From = m.From
	e.To = []string{msg.To}
	e.Subject = msg.Subject
	e.HTML = []byte(msg.HTMLContent)

	m.emailQueue <- e

	sendEmail := func() {
		for e := range m.emailQueue {
			err := m.pool.Send(e, m.timeout)
			if err != nil {
				m.errChan <- err
			}
		}
	}

	select {
		case <-ctx.Done():
			return
		
		default:
			go sendEmail() 
	}
}


func (m *Mailer) listen() {
	for {
		select{
		case <- m.doneChan:
			return

		case err := <-m.errChan:
			if err != nil {
				log.Error().Msg("Error on closing database" + err.Error())
			}
		}
	}
}