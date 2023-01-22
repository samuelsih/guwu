package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/samuelsih/guwu/business"
	"github.com/samuelsih/guwu/model"
	"github.com/samuelsih/guwu/pkg/errs"
	"github.com/samuelsih/guwu/pkg/securer"
)

const SESS_MAX_AGE = 60 * 60 * 24

type Deps struct {
	DB *sqlx.DB

	CreateSession  func(ctx context.Context, key string, in any, time int64) error
	DestroySession func(ctx context.Context, sessionID string) error
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	business.CommonResponse
	User model.User `json:"user"`
}

func (d *Deps) Login(ctx context.Context, in LoginInput, commonIn business.CommonInput) LoginOutput {
	var out LoginOutput

	if err := validEmail(in.Email); err != nil {
		out.RawError(400, err.Error())
		return out
	}

	if in.Password == "" {
		out.RawError(400, "password is required")
		return out
	}

	user, err := model.FindUserByEmail(ctx, d.DB, in.Email)
	if err != nil {
		out.SetError(err)
		return out
	}

	if !model.CheckUserPassword(user.Password.String, in.Password) {
		out.RawError(errs.KindBadRequest, "invalid credentials")
		return out
	}

	sessionID := xid.New().String()

	err = d.CreateSession(ctx, sessionID, user, int64(SESS_MAX_AGE))
	if err != nil {
		out.SetError(err)
		return out
	}

	encryptedSessionID, err := securer.Encrypt([]byte(sessionID))
	if err != nil {
		out.SetError(err)
		return out
	}

	out.User = user
	out.SessionID = encryptedSessionID
	out.SessionMaxAge = SESS_MAX_AGE
	out.SetOK()

	return out
}

type RegisterInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterOutput struct {
	business.CommonResponse
}

func (d *Deps) Register(ctx context.Context, in RegisterInput, commonIn business.CommonInput) RegisterOutput {
	var out RegisterOutput

	if err := validAccount(in.Username, in.Email, in.Password); err != nil {
		out.RawError(400, err.Error())
		return out
	}

	hashedPassword, err := model.HashPassword(in.Password)
	if err != nil {
		out.SetError(err)
		return out
	}

	_, err = model.InsertUser(ctx, d.DB, in.Username, in.Email, hashedPassword)
	if err != nil {
		out.SetError(err)
		return out
	}

	out.SetOK()
	return out
}

type LogoutOutput struct {
	business.CommonResponse
}

func (d *Deps) Logout(ctx context.Context, in business.CommonInput) LogoutOutput {
	var out LogoutOutput

	if in.SessionID == "" {
		out.RawError(400, "session id is required")
		return out
	}

	sessID, err := securer.Decrypt(in.SessionID)
	if err != nil {
		out.SetError(err)
		return out
	}

	err = d.DestroySession(ctx, string(sessID))

	if err != nil {
		out.SetError(err)
		return out
	}

	out.SessionID = ""
	out.SessionMaxAge = -1
	out.SetOK()

	return out
}
