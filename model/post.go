package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

type Post struct {
	ID          string `db:"id" json:"id"`
	UserID string `db:"user_id" json:"user_id"`
	Description string `db:"description" json:"description"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt NullTime `db:"updated_at" json:"updated_at,omitempty"`
}

type PostDeps struct {
	DB *sqlx.DB
}

func (p *PostDeps) Insert(ctx context.Context, description, userID string) (Post, statusCode, error) {
	query := `INSERT INTO posts(id, description, user_id) VALUES ($1, $2, $3)`

	post := Post{
		ID:  xid.New().String(),
		Description: description,
	}

	_, err := p.DB.ExecContext(ctx, query, post.ID, post.Description, userID)
	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "posts.InsertUser.ExecContext")
		return Post{}, BadRequest, err
	}

	return post, OK, nil
}

func (p *PostDeps) GetTimeline(ctx context.Context) ([]Post, int, error) {
	query := `SELECT p.id, p.description, p.created_at, p.updated_at FROM posts AS p JOIN users AS u on p.user_id = u.id`

	var posts []Post
	rows, err := p.DB.QueryxContext(ctx, query)

	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "posts.GetTimeline.QueryxContext")
		return nil, InternalServerError, err
	}

	defer rows.Close()

	for rows.Next() {
		var post Post
		err = rows.StructScan(&post)

		if err != nil {
			log.Debug().Stack().Err(err).Str("place", "posts.GetTimeline.StructScan")
			return nil, InternalServerError, err
		}

		posts = append(posts, post)
	}

	if len(posts) == 0 {
		return posts, NoContent, nil
	}

	return posts, OK, nil
}

func (p *PostDeps) Update(ctx context.Context, description string, postID, userID string) (Post, statusCode, error) {
	query := `UPDATE posts SET description = $1, updated_at = $2 WHERE user_id = $3 RETURNING *`

	var post Post
	err := p.DB.QueryRowxContext(ctx, query, description, time.Now(), userID).StructScan(&post)

	if err != nil {
		if err == sql.ErrNoRows {
			return post, BadRequest, errors.New(`unknown user`)
		}
		
		return post, InternalServerError, err
	}
	
	return post, OK, nil 
}