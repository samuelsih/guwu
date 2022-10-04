package model

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
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
		if errSQL, ok := err.(*pgconn.PgError); ok {
			switch errSQL.Code {
				case pgerrcode.ForeignKeyViolation:	
					return Post{}, errors.New(`unknown user`)
			}
		}
		log.Debug().Stack().Err(err).Str("place", "posts.InsertUser.ExecContext")
		return Post{}, errors.New(`can't add posts right now`)
	}

	return post, nil
}

func (p *PostDeps) GetTimeline(ctx context.Context) ([]Post, error) {
	query := `SELECT id, description, created_at, updated_at FROM posts JOIN users on posts.user_id = users.id`

	var posts []Post
	rows, err := p.DB.QueryxContext(ctx, query)

	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "posts.GetTimeline.QueryxContext")
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var post Post
		err = rows.StructScan(&post)

		if err != nil {
			log.Debug().Stack().Err(err).Str("place", "posts.GetTimeline.StructScan")
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (p *PostDeps) Update(ctx context.Context, description string, userID string) (Post, error) {
	query := `UPDATE posts SET description = $1, updated_at = $2 WHERE user_id = $3 RETURNING *`

	var post Post
	err := p.DB.QueryRowxContext(ctx, query, description, time.Now(), userID).StructScan(&post)

	if err != nil {
		return post, err
	}
	
	return post, nil 
}