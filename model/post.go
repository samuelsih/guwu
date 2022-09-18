package model

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/rs/xid"
)

type Post struct {
	ID          string `db:"id" json:"id"`
	Description string `db:"description" json:"description"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

type PostDeps struct {
	DB *sqlx.DB
}

func (p *PostDeps) Insert(ctx context.Context, description, userID string) (Post, error) {
	query := `INSERT INTO post(id, description, user_id) VALUES ($1, $2, $3)`

	post := Post{
		ID:  xid.New().String(),
		Description: description,
	}

	_, err := p.DB.ExecContext(ctx, query, post.ID, post.Description, userID)
	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "posts.InsertUser.ExecContext")
		return Post{}, wrapErr(err)
	}

	return post, nil
}

func (p *PostDeps) GetUserAllPosts(ctx context.Context, userID string) ([]Post, error) {
	query := `SELECT id, description, created_at FROM posts WHERE user_id = $1`

	var posts []Post
	rows, err := p.DB.QueryxContext(ctx, query, userID)

	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "posts.GetUserAllPosts.QueryxContext")
		return nil, errors.New(`no posts on this user`)
	}

	for rows.Next() {
		var post Post
		err = rows.StructScan(&post)

		if err != nil {
			log.Debug().Stack().Err(err).Str("place", "posts.GetUserAllPosts.StructScan")
			return nil, errors.New(`error getting posts on this user`)
		}

		posts = append(posts, post)
	}

	if posts == nil {
		return nil, errors.New(`no posts found`)
	}

	return posts, nil
}

func wrapErr(errSQL error) error {
	if err, ok := errSQL.(*pgconn.PgError); ok {
		switch err.Code {
			case pgerrcode.ForeignKeyViolation:	
				return errors.New(`unknown user`)
		}
	}

	return errors.New(`internal server error, please try again later`)
}
