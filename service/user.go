package service

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/model"
)

type User struct {
	DB *sqlx.DB
}

type UserFindOut struct {
	CommonResponse
	Users []model.User `json:"users,omitempty"`
}

func (u *User) FindUser(ctx context.Context, username string) UserFindOut {
	var out UserFindOut

	user := model.UserDeps{DB: u.DB}

	result, err := user.FindUserByUsername(ctx, username)
	if err != nil {
		out.SetError(http.StatusBadRequest, err.Error())
		return out
	}

	if len(result) == 0 {
		out.SetError(http.StatusBadRequest, "No users found")
		return out
	}

	out.Users = result
	out.SetOK()
	return out
}
