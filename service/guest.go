package service

import (
	"context"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/samuelsih/guwu/model"
)

type Guest struct {
	DB *sqlx.DB
	SessionDB *redis.Client
}

type GuestLoginIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GuestLoginOut struct {
	CommonResponse
	User *model.User `json:"user,omitempty"`
}

func (u *Guest) Login(ctx context.Context, in *GuestLoginIn) GuestLoginOut {
	var out GuestLoginOut

	user := model.UserDeps{DB: u.DB}

	result, err := user.GetUserByEmail(in.Email)
	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "user.GetUserByEmail")
		out.SetError(http.StatusBadRequest, err.Error())
		return out
	}

	if !result.PasswordMatches(in.Password) {
		log.Debug().Stack().Err(err).Str("place", "user.PasswordMatches")
		out.SetError(http.StatusBadRequest, `User or password does not match`)
		return out
	}

	out.SetOK()
	out.User = result
	return out
}

type GuestRegisterIn struct {
	Email    string `json:"email"`
	Username string `json:"name"`
	Password string `json:"password"`
}

type GuestRegisterOut struct {
	CommonResponse
	User *model.User `json:"user,omitempty"`
}

func (u *Guest) Register(ctx context.Context, in *GuestRegisterIn) GuestRegisterOut {
	var out GuestRegisterOut

	user := model.UserDeps{DB: u.DB}

	if err := validateSignIn(*in); err != nil {
		log.Debug().Stack().Err(err).Str("place", "validateSignIn")
		out.SetError(http.StatusBadRequest, err.Error())
		return out
	}

	result, err := user.Insert(ctx, in.Username, in.Email, in.Password)
	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "user.Insert")
		out.SetError(http.StatusBadRequest, err.Error())
		return out
	}

	out.User = result
	out.SetOK()
	return out
}

type GuestLogoutOut struct {
	CommonResponse
}

func (u *Guest) Logout(ctx context.Context, sessionID string) GuestLogoutOut {
	var out GuestLogoutOut

	sess := model.SessionDeps{Conn: u.SessionDB}

	err := sess.Delete(ctx, sessionID)
	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "session.Delete")
		out.SetError(http.StatusBadRequest, err.Error())
		return out
	}

	out.SetOK()
	return out
}