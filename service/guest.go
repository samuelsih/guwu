package service

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/model"
)

type Guest struct {
	DB *sqlx.DB
}

type GuestLoginIn struct {
	UserCred string `json:"user_cred"`
	Password string `json:"password"`
}

type GuestLoginOut struct {
	CommonResponse
	User *model.User
}

func (u *Guest) Login(ctx context.Context, in *GuestLoginIn) GuestLoginOut {
	var out GuestLoginOut

	user := model.NewUser(u.DB)
	
	user, err := user.GetUserByEmailOrUsername(in.UserCred)
	if err != nil {
		out.SetError(http.StatusBadRequest, err.Error())
		return out
	}

	if !user.PasswordMatches(in.Password) {
		out.SetError(http.StatusBadRequest, `User or password does not match`)
		return out
	} 

	out.SetOK()
	out.User = user.Clean()
	return out
}

type GuestRegisterIn struct {
	Email string `json:"email"`
	Name string `json:"name"`
	Password string `json:"password"`
}

type GuestRegisterOut struct {
	CommonResponse
	User *model.User `json:"user,omitempty"`
}


func (u *Guest) Register(ctx context.Context, in *GuestRegisterIn) GuestRegisterOut {
	var out GuestRegisterOut

	user := model.NewUser(u.DB)
	
	if user.EmailExists(ctx, in.Email) {
		out.SetError(http.StatusBadRequest, `Email already exists`)
		return out
	}

	user.Email = in.Email
	user.Name = in.Name
	user.Password = in.Password

	user, err := user.Insert(ctx)
	if err != nil {
		out.SetError(http.StatusInternalServerError, "Cant insert user to database:" + err.Error())
		return out
	}

	out.User = user.Clean()
	return out
}
