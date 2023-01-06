package auth

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/samuelsih/guwu/business"
	"github.com/samuelsih/guwu/model"
	"github.com/samuelsih/guwu/pkg/redis"
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
		out.SetError(400, err.Error())
		return out
	}

	if in.Password == "" {
		out.SetError(400, errPasswordRequired.Error())
		return out
	}

	user := model.NewUser(d.DB)

	if !user.FindByEmail(ctx, in.Email) {
		out.SetError(user.StatusCode, user.Error())
		return out
	}

	if !user.CheckPassword(in.Password) {
		out.SetError(user.StatusCode, user.Error())
		return out
	}

	sessionID := xid.New().String()
	authenticatedUser := user.Cleanup()

	err := d.CreateSession(ctx, sessionID, authenticatedUser, int64(SESS_MAX_AGE))

	if err != nil {
		out.SetError(500, err.Error())
		return out
	}

	encryptedSessionID, err := securer.Encrypt([]byte(sessionID))	
	if err != nil {
		out.SetError(500, err.Error())
		return out
	}

	out.User = authenticatedUser
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
		out.SetError(400, err.Error())
		return out
	}

	user := model.NewUser(d.DB)

	user.SetPassword(in.Password)

	if user.Err != nil {
		out.SetError(user.StatusCode, user.Error())
		return out
	}

	user.Email = in.Email
	user.Username = in.Username

	if !user.Insert(ctx) {
		out.SetError(user.StatusCode, user.Error())
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
		out.SetError(400, "session id is required")
		return out
	}

	sessID, err := securer.Decrypt(in.SessionID)
	if err != nil {
		out.SetError(400, err.Error())
		return out
	}

	err =  d.DestroySession(ctx, string(sessID))

	if err != nil {
		if errors.Is(err, redis.ErrUnknownKey) {
			out.SetError(400, `unknown session`)
			return out
		}

		out.SetError(500, err.Error())
		return out
	}

	out.SessionID = ""
	out.SessionMaxAge = -1
	out.SetOK()

	return out
}
