package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/pkg/errs"
	"github.com/samuelsih/guwu/pkg/pgerr"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string     `db:"id" json:"id"`
	Username  string     `db:"username" json:"username"`
	Email     string     `db:"email" json:"email"`
	Password  NullString `db:"password" json:"-"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt NullTime   `db:"updated_at" json:"updated_at,omitempty"`
}

func FindUserByEmail(ctx context.Context, db *sqlx.DB, email string) (User, error) {
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = $1`
	const op = errs.Op("user.FindByEmail")
	var user User

	err := db.GetContext(ctx, &user, query, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errs.E(op, errs.KindBadRequest, err, "unknown user")
		}

		return user, errs.E(op, errs.KindUnexpected, err, "cannot get user email")
	}

	if user.ID == "" {
		return user, errs.E(op, errs.KindUnexpected, errors.New("user id is nil"), "cannot get user email")
	}

	return user, nil
}

func CheckUserPassword(userPassword, incomingPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(incomingPassword))
	return err == nil
}

func InsertUser(ctx context.Context, db *sqlx.DB, username, email, password string) (User, error) {
	query := `
		INSERT INTO users(username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, created_at;
	`
	const op = errs.Op("user.Insert")
	var user User

	err := db.QueryRowContext(ctx, query, username, email, password).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		if column, e := pgerr.UniqueColumn(err); e != nil {
			clientMsg := fmt.Sprintf("%v already taken, please take another %v", column, column)
			return user, errs.E(op, errs.KindBadRequest, e, clientMsg)
		}

		return user, errs.E(op, errs.KindUnexpected, err, "unexpected error.")
	}

	return user, nil
}

func HashPassword(password string) (string, error) {
	const op = errs.Op("user.HashPassword")
	
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errs.E(op, errs.KindUnexpected, err, "unexpected error.")
	}

	return string(hashed), nil
}
