package mail

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer *gomail.Dialer
	mailChan chan *gomail.Message
	sendCloser gomail.SendCloser

	limiter chan struct{}
	doneChan chan struct{}
}

type Message struct {
	To string
	Subject string
	PlainContent string
	HTMLContent string
}

func NewMailer(host string, port int, user, password string) *Mailer {
	dialer := gomail.NewDialer(host, port, user, password)
	mailer := &Mailer{
		dialer: dialer,
		mailChan: make(chan *gomail.Message, 10),
		doneChan: make(chan struct{}),
		limiter: make(chan struct{}, 10),
	}
	return mailer
}

func (m *Mailer) Send(msg Message) {
	message := gomail.NewMessage()
	message.SetHeader("From", "guwu@info.com")
	message.SetHeader("To", msg.To)
	message.SetHeader("subject", msg.Subject)
	message.SetBody("text/html", msg.HTMLContent)

	m.mailChan <- message
}

func (m *Mailer) Listen() {
	for {
		select{
		case <- m.doneChan:
			close(m.mailChan)
			close(m.doneChan)
			close(m.limiter)
			return

		case msg := <- m.mailChan:
			m.limiter <- struct{}{} // will blocking if full

			go func() {
				defer func(){
					<- m.limiter
				}()

				if err := gomail.Send(m.sendCloser, msg); err != nil {
					log.Error().Msg("Error on mailer: " + err.Error())
				}
			}()
		}
	}
}

func (m *Mailer) Close() { 
	m.doneChan <- struct{}{} 
}