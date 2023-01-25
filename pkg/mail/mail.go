package mail

import (
	"context"
	_ "embed"
	"errors"
	ht "html/template"
	tt "text/template"

	"github.com/samuelsih/guwu/pkg/errs"
	"github.com/wneessen/go-mail"
)

type MsgType int

var (
	//go:embed otp.html
	otpHTML string

	//go:embed otp.txt
	otpTxt string
)

const (
	OTPMsg MsgType = iota
	RecoverPasswdMsg
)

type Client struct {
	client      *mail.Client
	senderName  string
	senderEmail string
}

type Param struct {
	Name          string
	Email         string
	Subject       string
	TemplateTypes MsgType
}

type OTPTplData struct {
	Username string
	OTP      string
}

type RecoverPasswdTplData struct {
	Username      string
	GeneratedLink string
}

func NewClient(host string, port int, username, password, senderName, senderEmail string) (Client, error) {
	const op = errs.Op("mail.NewClient")

	client, err := mail.NewClient(host,
		mail.WithPort(port),
		mail.WithTLSPolicy(mail.TLSOpportunistic),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username),
		mail.WithPassword(password),
	)

	if err != nil {
		return Client{}, errs.E(op, errs.KindUnexpected, err, "unexpected server error")
	}

	return Client{client: client, senderName: senderName, senderEmail: senderEmail}, nil
}

func (c Client) Send(ctx context.Context, param Param, tplData any) error {
	const op = errs.Op("mail.Send")

	htpl, txtpl, err := getTemplateFromType(param.TemplateTypes)
	if err != nil {
		return errs.E(op, errs.GetKind(err), err, "cannot generate message")
	}

	m := mail.NewMsg()

	if err := m.FromFormat(c.senderName, c.senderEmail); err != nil {
		return errs.E(op, errs.KindUnexpected, err, "unexpected error")
	}

	if err := m.AddToFormat(param.Name, param.Email); err != nil {
		return errs.E(op, errs.KindUnexpected, err, "unexpected error to format")
	}

	m.Subject(param.Subject)
	m.SetMessageID()
	m.SetDate()

	if err := m.SetBodyHTMLTemplate(htpl, tplData); err != nil {
		return errs.E(op, errs.KindUnexpected, err, "can't set body email")
	}

	if err := m.AddAlternativeTextTemplate(txtpl, tplData); err != nil {
		return errs.E(op, errs.KindUnexpected, err, "can't set body text")
	}

	if err := c.client.DialAndSendWithContext(ctx, m); err != nil {
		return errs.E(op, errs.KindUnexpected, err, "cant send message")
	}

	return nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func getTemplateFromType(msgType MsgType) (*ht.Template, *tt.Template, error) {
	const op = errs.Op("mail.getTemplateFromType")

	switch msgType {
	case OTPMsg:
		htpl, err := ht.New("htmltpl").Parse(otpHTML)
		if err != nil {
			return nil, nil, errs.E(op, errs.KindUnexpected, err, "unexpected error generating html")
		}

		txtpl, err := tt.New("texttpl").Parse(otpTxt)
		if err != nil {
			return nil, nil, errs.E(op, errs.KindUnexpected, err, "unexpected error generating txt")
		}

		return htpl, txtpl, nil

	default:
		return nil, nil, errs.E(op, errs.KindUnexpected, errors.New("unknown template"), "unexpected error generating message")
	}
}
