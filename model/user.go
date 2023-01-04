package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/pkg/pgerr"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	db         *sqlx.DB `json:"-"`
	Err        error    `json:"-"`
	StatusCode int      `json:"-"`

	ID        string     `db:"id" json:"id"`
	Username  string     `db:"username" json:"username"`
	Email     string     `db:"email" json:"email"`
	Password  NullString `db:"password" json:"-"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt NullTime   `db:"updated_at" json:"updated_at,omitempty"`
}

func NewUser(db *sqlx.DB) *User {
	return &User{db: db}
}

func (u *User) FindByEmail(ctx context.Context, email string) bool {
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = $1`
	err := u.db.GetContext(ctx, u, query, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			u.Err = fmt.Errorf("unknown user")
			u.StatusCode = 400
			return false
		}

		log.Printf("User.FindByEmail: %v", err)
		u.Err = fmt.Errorf("can't find user's email, try again later")
		u.StatusCode = 500
		return false
	}

	return (u.ID != "") && (err == nil)
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password.String), []byte(password))

	if err != nil {
		u.Err = fmt.Errorf(`unknown type of password`)
		u.StatusCode = 400
		return false
	}

	return true
}

func (u *User) Insert(ctx context.Context) bool {
	query := `
		INSERT INTO users(username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at;
	`

	err := u.db.QueryRowContext(ctx, query, u.Username, u.Email, u.Password.String).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		if column, e := pgerr.UniqueColumn(err); e != nil {
			u.Err = fmt.Errorf("%v already taken, please take another %v", column, column)
			u.StatusCode = 400
			return false
		}

		log.Printf("User.Insert: %v", err)
		u.Err = fmt.Errorf(`can't insert this user to database`)
		u.StatusCode = 500
		return false
	}

	return true
}

func (u *User) SetPassword(password string) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.Err = err
		u.StatusCode = 500
	}
	u.Password.String = string(hashed)
}

func (u User) Cleanup() User {
	u.db = nil
	u.Err = nil
	return u
}

func (u *User) Error() string {
	return u.Err.Error()
}
