package model

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("User not found")
)

type User struct {
	db *sqlx.DB `json:"-"`
	ID        string `db:"id"`
	Name      string `db:"name"`
	Email     string `db:"email"`
	Password  string `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

func NewUser(db *sqlx.DB) *User {
	return &User{db: db}
}

func (u *User) PasswordMatches(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) Insert(ctx context.Context) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u.ID = xid.New().String()
	u.Password = string(hashedPassword)
	u.CreatedAt = time.Now()

	query := `INSERT INTO users (id, name, email, password, created_at) 
			VALUES ($1, $2, $3, $4, $5, $6)`	

	result, err := u.db.ExecContext(ctx, query, u.ID, u.Name, u.Email, u.Password, u.CreatedAt)
	if err != nil {
		return nil, err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) GetUserByEmailOrUsername(param string) (*User, error) {
	var user User

	query := `SELECT id, name, username, email FROM users WHERE email = $1 or username = $2`

	err := u.db.Get(&u, query, param, param)
	
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}

	return &user, err
}

func (u *User) EmailExists(ctx context.Context, email string) bool {
	query := `SELECT email FROM users WHERE email = $1`

	err := u.db.GetContext(ctx, u, query, email)
	if err != nil {
		log.Printf("Error in User.EmailExists: %v", err)
		return false
	}

	return err == nil && u.ID != ""
}

func (u *User) Clean() *User {
	u.db = nil
	return u
}