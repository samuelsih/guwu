package service

import "github.com/jmoiron/sqlx"

type User struct {
	DB *sqlx.DB
}

type UserGetProfileIn struct {
	CommonRequest
}

type UserGetProfileOut struct {
	CommonResponse
}