package mailer

import (
	"gopkg.in/gomail.v2"
)

type mailer struct {
	email string
	password string
	messagePool chan string
	senderPool chan gomail.SendCloser
}

func NewMailer(email, password string) *mailer {
	mail := &mailer{
		email: email,
		password: password,
		messagePool: make(chan string),
		senderPool: make(chan gomail.SendCloser),
	}

	go mail.listenForMail()

	return mail
}

func (m *mailer) listenForMail() {

}

func (m *mailer) Send() error {
	return nil
}

func (m *mailer) Close() error {
	return nil
}