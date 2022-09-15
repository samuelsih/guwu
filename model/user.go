package model

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

type UserDeps struct {
	DB *sqlx.DB
}

func (u *UserDeps) Insert(ctx context.Context, username, email, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := User{
		ID: xid.New().String(),
		Username: username,
		Email: email,
		Password: string(hashedPassword),
		CreatedAt: time.Now().UTC(),
	}

	query := `INSERT INTO users (id, username, email, password, created_at) 
			VALUES ($1, $2, $3, $4, $5)`

	_, err = u.DB.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		return nil, wrapErr(err, "User")
	}

	return &user, nil
}

func (u *UserDeps) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, username, email, password FROM users WHERE email = $1`

	var user User
	err := u.DB.Get(&user, query, email)

	if err != nil {
		return nil, wrapErr(err, "Email")
	}

	return &user, nil
}

func (u *UserDeps) FindUserByUsername(ctx context.Context, username string) ([]*User, error) {
	query := `SELECT username FROM users WHERE username ILIKE $1`

	var users []*User
	err := u.DB.SelectContext(ctx, &users, query, username)

	if err != nil {
		return nil, wrapErr(err, "Username")
	}

	return users, nil
}

func (u *User) PasswordMatches(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
