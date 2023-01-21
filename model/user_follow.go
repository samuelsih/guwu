package model

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/pkg/pgerr"
)

type UserFollow struct {
	db         *sqlx.DB `json:"-"`
	Err        error    `json:"-"`
	StatusCode int      `json:"-"`

	UserID string `json:"user_id"`
	UserFollowID string `json:"user_follow_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func NewUserFollow(db *sqlx.DB) *UserFollow {
	return &UserFollow{db: db}
}

func (u *UserFollow) Insert(ctx context.Context, userID string, userWantsToFollow string) bool {
	q := `INSERT INTO user_follows (user_id, user_follow_id) VALUES ($1, $2)`

	_, err := u.db.ExecContext(ctx, q, userID, userWantsToFollow)
	if err != nil {
		if column, e := pgerr.ForeignKeyColumn(err); e != nil {
			u.Err = fmt.Errorf("unknown id for %v", column)
			u.StatusCode = 400
			return false
		}

		u.Err = fmt.Errorf("internal error"); log.Println("error db:", err)
		u.StatusCode = 500
		return false
	}

	return true
}

func (u *UserFollow) Delete(ctx context.Context, userID string, userWantsToUnfollow string) bool {
	q := `DELETE FROM user_follow WHERE user_id = $1 AND user_follow_id = $2`

	_, err := u.db.ExecContext(ctx, q, userID, userWantsToUnfollow)
	if err != nil {
		if column, e := pgerr.ForeignKeyColumn(err); e != nil {
			u.Err = fmt.Errorf("unknown id for %v", column)
			u.StatusCode = 400
			return false
		}

		u.Err = fmt.Errorf("internal error")
		u.StatusCode = 500
		return false
	}

	return true
}