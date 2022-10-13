package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"html/template"
	"net/http"

	"github.com/samuelsih/guwu/mail"
	"golang.org/x/crypto/bcrypt"
)

type Verification struct {
	SendEmail func(ctx context.Context, msg mail.Message)
}

type VerifOut struct {
	CommonResponse
}

func (v *Verification) Send(ctx context.Context, email string) VerifOut {
	var out VerifOut

	htmlContent, err := buildHTMLMessage(email)
	if err != nil {
		out.SetError(http.StatusInternalServerError, err.Error())
		return out
	}

	msg := mail.Message{
		To: email,
		Subject: "Verification",
		PlainContent: "ini verif: " + "http://localhost:8080/verif?data=" + generateToken(email),
		HTMLContent: htmlContent,
	}

	v.SendEmail(ctx, msg)

	out.SetOK()
	return out
}

func generateToken(email string) string {
    hash, _ := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
    return base64.StdEncoding.EncodeToString(hash)
}

func buildHTMLMessage(email string) (string, error) {
	templateToRender := "../mail/verif.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "link", "http://localhost:8080/verif?data=" + email); err != nil {
		return "", err
	}

	return tpl.String(), nil
}