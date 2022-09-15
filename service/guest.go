package service

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/samuelsih/guwu/model"
)

type Guest struct {
	DB *sqlx.DB
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

	user := model.NewUser(u.DB)

	err := user.GetUserByEmail(in.Email)
	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "user.GetUserByEmail")
		out.SetError(http.StatusBadRequest, err.Error())
		return out
	}

	if !user.PasswordMatches(in.Password) {
		log.Debug().Stack().Err(err).Str("place", "user.PasswordMatches")
		out.SetError(http.StatusBadRequest, `User or password does not match`)
		return out
	}

	out.SetOK()
	out.User = user
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

	user := model.NewUser(u.DB)

	if err := validateSignIn(*in); err != nil {
		log.Debug().Stack().Err(err).Str("place", "validateSignIn")
		out.SetError(http.StatusBadRequest, err.Error())
		return out
	}

	user.Email = in.Email
	user.Username = in.Username
	user.Password = in.Password

	user, err := user.Insert(ctx)
	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "user.Insert")
		out.SetError(http.StatusBadRequest, err.Error())
		return out
	}

	out.User = user
	out.SetOK()
	return out
}
