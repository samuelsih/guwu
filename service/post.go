package service

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/model"
)

type PostDeps struct {
	DB *sqlx.DB
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
		out.SetError(http.StatusBadRequest, err.Error())
		return out
	}

	out.Post = result
	out.SetOK()
	return out
}


type PostGetUserAllIn struct {
	CommonRequest
}

type PostGetUserAllOut struct {
	CommonResponse
	Posts []model.Post
}

func (p *PostDeps) GetAllByUser(ctx context.Context, in *PostGetUserAllIn, userSess model.Session) PostGetUserAllOut {
	var out PostGetUserAllOut

	post := model.PostDeps{DB: p.DB}
	
	result, err := post.GetUserAllPosts(ctx, in.UserSession.ID)
	if err != nil {
		out.SetError(http.StatusBadRequest, err.Error())
		return out
	}

	out.Posts = result
	out.SetOK()
	return out
}