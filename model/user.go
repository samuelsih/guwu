package model

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)


type User struct {
	db *sqlx.DB `json:"-"`
	ID        string `db:"id" json:"id"`
	Username      string `db:"username" json:"username"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type UserDeps struct {
	DB *sqlx.DB
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

	query := `INSERT INTO users (id, username, email, password, created_at) 
			VALUES ($1, $2, $3, $4, $5)`	

	_, err = u.db.ExecContext(ctx, query, u.ID, u.Username, u.Email, u.Password, u.CreatedAt)
	if err != nil {
		return nil, wrapErr(err, "User")
	}

	return u.clean(), nil
}

func (u *User) GetUserByEmail(email string) error {
	query := `SELECT id, username, email, password FROM users WHERE email = $1`

	err := u.db.Get(u, query, email)
	
	if err != nil {
		println("Error in getting user: ", err.Error())
		return wrapErr(err, "Email")
	}

	return err
}

func (u *User) FindUserByUsername(ctx context.Context, username string) ([]*User, error) {
	query := `SELECT username FROM users WHERE username ILIKE $1`

	var users []*User
	err := u.db.SelectContext(ctx, &users, query, username)

	if err != nil {
		return nil, wrapErr(err, "Username")
	}

	return users, nil
} 

func (u *User) clean() *User {
	u.db = nil
	return u
}

