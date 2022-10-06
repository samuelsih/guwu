package service

import (
	"context"
	"net/http"
	"unicode/utf8"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/samuelsih/guwu/model"
)

type Post struct {
	DB *sqlx.DB
}

type PostTimelineOut struct {
	CommonResponse
	Posts []model.Post `json:"posts,omitempty"`
}

func (p *Post) Timeline(ctx context.Context) PostTimelineOut {
	var out PostTimelineOut

	post := model.PostDeps{DB: p.DB}

	result, statusCode, err := post.GetTimeline(ctx)
	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "posts.GetTimeline")
		out.SetError(statusCode, err.Error())
		return out
	}

	out.Posts = result
	out.SetOK()
	return out
}

type PostInsertIn struct {
	CommonRequest
	Description string `json:"description"`
}

type PostInsertOut struct {
	CommonResponse
	Post model.Post `json:"post"` 
}

func (p *Post) Insert(ctx context.Context, in *PostInsertIn) PostInsertOut {
	var out PostInsertOut

	if in.Description == "" {
		out.SetError(http.StatusBadRequest, `description is required`)
		return out
	}

	if utf8.RuneCountInString(in.Description) > 512 {
		out.SetError(http.StatusBadRequest, `description is too long, max is 512 characters`)
		return out
	}

	post := model.PostDeps{DB: p.DB}

	result, statusCode, err := post.Insert(ctx, in.Description, in.UserSession.ID)
	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "posts.Insert")
		out.SetError(statusCode, err.Error())
		return out
	}

	out.Post = result
	out.SetOK()
	return out
}

type PostEditIn struct {
	CommonRequest
	Description string `json:"description"`
}

type PostEditOut struct {
	CommonResponse
	Post model.Post `json:"post"`
}

func (p *Post) Edit(ctx context.Context, postID string, in *PostEditIn) PostEditOut {
	var out PostEditOut

	post := model.PostDeps{DB: p.DB}

	result, statusCode, err := post.Update(ctx, in.Description, postID, in.UserSession.ID)
	if err != nil {
		out.SetError(statusCode, err.Error())
		return out
	}

	if result == (model.Post{}) {
		out.SetError(http.StatusInternalServerError, `internal server error`)
		return out
	}

	out.Post = result
	out.SetOK()
	return out
}