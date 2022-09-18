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
	SessionID string `json:"-"`
	User model.User `json:"user,omitempty"`
}

func (u *Guest) Login(ctx context.Context, in *GuestLoginIn) GuestLoginOut {	
	var out GuestLoginOut

	if len(in.Email) < 3 {
		out.SetError(400, `invalid email`)
		return out
	}
	if len(in.Password) < 3 {
		out.SetError(400, `invalid password`)
		return out
	}

	user := model.UserDeps{DB: u.DB}

	result, err := user.GetUserByEmail(in.Email)
	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "user.GetUserByEmail")
		out.SetError(http.StatusBadRequest, err.Error())
		return out
	}

	if !result.PasswordMatches(in.Password) {
		log.Debug().Stack().Err(err).Str("place", "user.PasswordMatches")
		out.SetError(http.StatusBadRequest, `email or password does not match`)
		return out
	}

	sess := model.SessionDeps{Conn: u.SessionDB}
	sessionData := model.Session{
		ID: result.ID, 
		Username: result.Username,
		Email: result.Email,
	}

	sessionID, err := sess.Save(ctx, sessionData)
	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "user.SaveSession")
		out.SetError(http.StatusInternalServerError, `error creating user session`)
		return out
	}

	out.User = result
	out.SessionID = sessionID
	out.SetOK()
	return out
}

type GuestRegisterIn struct {
	Email    string `json:"email"`
	Username string `json:"name"`
	Password string `json:"password"`
}

type GuestRegisterOut struct {
	CommonResponse
	SessionID string `json:"-"`
	User model.User `json:"user,omitempty"`
}

func (u *Guest) Register(ctx context.Context, in *GuestRegisterIn) GuestRegisterOut {
	var out GuestRegisterOut

	user := model.UserDeps{DB: u.DB}

	if err := validateSignIn(*in); err != nil {
		log.Debug().Stack().Err(err).Str("place", "validateSignIn")
		out.SetError(http.StatusBadRequest, err.Error())
		return out
	}

	_, err := user.GetUserByEmail(in.Email)
	if err == nil {
		log.Debug().Stack().Err(err).Str("place", "validateSignIn")
		out.SetError(http.StatusBadRequest, `email already exists`)
		return out
	}

	result, err := user.Insert(ctx, in.Username, in.Email, in.Password)
	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "user.Insert")
		out.SetError(http.StatusInternalServerError, err.Error())
		return out
	}

	sess := model.SessionDeps{Conn: u.SessionDB}
	sessionData := model.Session{
		ID: result.ID, 
		Username: result.Username,
		Email: result.Email,
	}

	sessionID, err := sess.Save(ctx, sessionData)
	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "user.SaveSession")
		out.SetError(http.StatusInternalServerError, `error creating user session`)
		return out
	}

	out.User = result
	out.SessionID = sessionID
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