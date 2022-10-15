package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"path"

	"github.com/samuelsih/guwu/mail"
	"github.com/samuelsih/guwu/internal"
	"golang.org/x/crypto/bcrypt"
)

type Verification struct {
	SendEmail func(ctx context.Context, msg mail.Message)
}

type VerifSendOut struct {
	CommonResponse
}

func (v *Verification) Send(ctx context.Context, email string) VerifSendOut {
	var out VerifSendOut

	if err := validateEmail(email); err != nil {
		out.SetError(http.StatusBadRequest, err.Error())
		return out
	}

	ok, err := internal.IsSpamEmail(ctx, email)
	if err != nil {
		out.SetError(http.StatusInternalServerError, err.Error())
		return out
	}

	if !ok {
		out.SetError(http.StatusBadRequest, "this email is considered spam.")
		return out
	}

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

type VerifyCheckOut struct {
	Token string `json:"token"`
	CommonResponse
}

func (v *Verification) Check(ctx context.Context, link string) VerifyCheckOut {
	var out VerifyCheckOut

	token, err := parseLink(link)
	if err != nil {
		out.SetError(http.StatusBadRequest, err.Error())
		return out
	}

	out.Token = token
	out.SetOK()
	return out
} 

func generateToken(email string) string {
    hash, _ := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
    return base64.StdEncoding.EncodeToString(hash)
}

func buildHTMLMessage(email string) (string, error) {
	templateToRender := "./mail/verif.gohtml"

	t, err := template.ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, map[string]string{
		"link": "http://localhost:8080/verif?data=" + generateToken(email),
	}); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func parseLink(link string) (string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return "", errors.New("not a valid url")
	}

	return path.Base(u.Path), nil
}