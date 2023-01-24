package mail

import (
	"context"
	"time"

	"github.com/samuelsih/guwu/pkg/errs"
	"github.com/wneessen/go-mail"
)

type Client struct {
	client *mail.Client
	// limiter soon
}

func NewClient(host string, port int, username, password string, timeout time.Duration) (Client, error) {
	const op = errs.Op("mail.NewClient")

	client, err := mail.NewClient(host,
		mail.WithPort(port),
		mail.WithTLSPolicy(mail.TLSOpportunistic),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTimeout(timeout),
	)

	if err != nil {
		return Client{}, errs.E(op, errs.KindUnexpected, err, "unexpected server error")
	}

	return Client{client: client}, nil
}

func (c Client) Send(ctx context.Context, from, to, subject, body string) error {
	const op = errs.Op("mail.Send")

	mailer := mail.NewMsg()
	mailer.Subject(subject)
	mailer.SetBodyString(mail.TypeTextPlain, body)

	if err := mailer.From(from); err != nil {
		return errs.E(op, errs.KindUnexpected, err, "cant send email from this user")
	}

	if err := mailer.To(to); err != nil {
		return errs.E(op, errs.KindUnexpected, err, "cant send email to this user")
	}

	if err := c.client.DialAndSendWithContext(ctx, mailer); err != nil {
		return errs.E(op, errs.KindUnexpected, err, "cant send message")
	}

	return nil
}

func (c *Client) Close() error {
	return c.client.Close()
}