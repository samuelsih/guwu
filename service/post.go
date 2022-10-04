package service

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/samuelsih/guwu/model"
)

type PostDeps struct {
	DB *sqlx.DB
}

type PostTimelineOut struct {
	CommonResponse
	Posts []model.Post `json:"posts,omitempty"`
}

func (p *PostDeps) Timeline(ctx context.Context) PostTimelineOut {
	var out PostTimelineOut

	post := model.PostDeps{DB: p.DB}

	result, err := post.GetTimeline(ctx)
	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "posts.GetTimeline")
		out.SetError(http.StatusBadRequest, err.Error())
		return out
	}

	if len(result) == 0 {
		log.Debug().Stack().Err(err).Str("place", "posts.len(result")
		out.SetError(http.StatusNoContent, `no posts for now`)
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

func (p *PostDeps) Insert(ctx context.Context, in *PostInsertIn) PostInsertOut {
	var out PostInsertOut

	post := model.PostDeps{DB: p.DB}

	result, err := post.Insert(ctx, in.Description, in.UserSession.ID)
	if err != nil {
		log.Debug().Stack().Err(err).Str("place", "posts.Insert")
		out.SetError(http.StatusBadRequest, err.Error())
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

func (p *PostDeps) Edit(ctx context.Context, in *PostEditIn) PostEditOut {
	var out PostEditOut

	post := model.PostDeps{DB: p.DB}

	result, err := post.Update(ctx, in.Description, in.UserSession.ID)
	if err != nil {
		out.SetError(http.StatusBadRequest, err.Error())
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