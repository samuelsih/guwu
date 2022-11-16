package model

import (
	"context"
	_ "embed"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

type User struct {
	ID        string    `db:"id" json:"id"`
	Name  		string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt NullTime `db:"updated_at" json:"updated_at,omitempty"`
}

type UserDeps struct {
	DB *sqlx.DB
}

func (u *UserDeps) Insert(ctx context.Context, name, email, password string) (User, statusCode, error) {
	user := User{
		ID: xid.New().String(),
		Name: name,
		Email: email,
		CreatedAt: time.Now().UTC(),
	}

	query := `INSERT INTO users (id, name, email, created_at) 
			VALUES ($1, $2, $3, $4, $5)`

	_, err := u.DB.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.CreatedAt)
	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "user.Insert.ExecContext")
		return User{}, InternalServerError, err
	}

	return user, Created, nil
}

func (u *UserDeps) GetUserByEmail(email string) (User, statusCode, error) {
	query := `SELECT id, name, email, password FROM users WHERE email = $1`

	var user User
	err := u.DB.Get(&user, query, email)

	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "user.GetUserByEmail")
		return user, BadRequest, errors.New("user not found")
	}
	
	return user, OK, nil
}

func (u *UserDeps) FindUserByName(ctx context.Context, name string) ([]User, statusCode, error) {
	query := `SELECT name FROM users WHERE name ILIKE $1`

	var users []User
	err := u.DB.SelectContext(ctx, &users, query, name)

	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "user.FindUserByName")
		return nil, InternalServerError, errors.New("user not found")
	}

	return users, OK, nil
}