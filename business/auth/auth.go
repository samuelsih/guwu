package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/business"
	"github.com/samuelsih/guwu/model"
)

type Deps struct {
	DB *sqlx.DB
}

type LoginInput struct {
	business.CommonRequest
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	business.CommonResponse
	User model.User `json:"user"`
}

func (d *Deps) Login(ctx context.Context, in LoginInput) LoginOutput {
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

	out.User = user.Cleanup()
	out.SetOK()

	return out
}

type RegisterInput struct {
	business.CommonRequest
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterOutput struct {
	business.CommonResponse
}

func (d *Deps) Register(ctx context.Context, in RegisterInput) RegisterOutput {
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
