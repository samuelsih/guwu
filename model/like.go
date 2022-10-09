package model

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

type Like struct {
	ID     string `db:"id" json:"id"`
	PostID string `db:"post_id" json:"post_id"`
	UserID string `db:"user_id" json:"user_id"`
}

type LikeDeps struct {
	DB *sqlx.DB
}

func(l *LikeDeps) Insert(ctx context.Context, postID, userID string) (Like, statusCode, error) {
	query := `INSERT INTO likes(id, post_id, user_id) VALUES ($1, $2, $3)`

	like := Like{
		ID:  xid.New().String(),
		PostID: postID,
		UserID: userID,
	}

	_, err := l.DB.ExecContext(ctx, query, like.ID, like.PostID, like.UserID)
	if err != nil {
		if errSQL, ok := err.(*pgconn.PgError); ok {
			switch errSQL.Code {
				case pgerrcode.ForeignKeyViolation:	
					return Like{}, BadRequest, errors.New(`unknown user/posts`)
			}
		}
		log.Debug().Stack().Err(err).Str("place", "likes.Insert.ExecContext")
		return Like{}, BadRequest, err
	}

	return like, 200, nil
}