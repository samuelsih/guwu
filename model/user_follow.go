package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/pkg/errs"
	"github.com/samuelsih/guwu/pkg/pgerr"
)

type UserFollow struct {
	UserID       string    `json:"user_id"`
	UserFollowID string    `json:"user_follow_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

func FollowUser(ctx context.Context, db *sqlx.DB, userID string, userWantsToFollow string) error {
	q := `INSERT INTO user_follows (user_id, user_follow_id) VALUES ($1, $2)`
	const op = errs.Op("user_follow.FollowUser")

	_, err := db.ExecContext(ctx, q, userID, userWantsToFollow)
	if err != nil {
		if column, e := pgerr.ForeignKeyColumn(err); e != nil {
			clientMsg := fmt.Sprintf("unknown id for %v", column)
			return errs.E(op, errs.KindBadRequest, err, clientMsg)
		}

		return errs.E(op, errs.KindUnexpected, err, "cannot follow the user.")
	}

	return nil
}

func UnfollowUser(ctx context.Context, db *sqlx.DB, userID string, userWantsToUnfollow string) error {
	q := `DELETE FROM user_follow WHERE user_id = $1 AND user_follow_id = $2`
	const op = errs.Op("user_follow.UnfollowUser")

	_, err := db.ExecContext(ctx, q, userID, userWantsToUnfollow)
	if err != nil {
		if _, e := pgerr.ForeignKeyColumn(err); e != nil {
			return errs.E(op, errs.KindBadRequest, err, "unknown user")
		}

		return errs.E(op, errs.KindUnexpected, err, "cannot unfollow this user")
	}

	return nil
}
