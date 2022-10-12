package service

import (
	"context"
	"encoding/base64"

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

	msg := mail.Message{
		To: email,
		Subject: "Verification",
		PlainContent: "ini verif: " + "http://localhost:8080/verif?data=" + generateToken(email),
		HTMLContent: "</p>ini verif</p>",
	}

	v.SendEmail(ctx, msg)

	out.SetOK()
	return out
}

func generateToken(email string) string {
    hash, _ := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
    return base64.StdEncoding.EncodeToString(hash)
}
