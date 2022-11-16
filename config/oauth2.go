package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var GithubConfig = &oauth2.Config{
	ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	Endpoint:     github.Endpoint,
	Scopes:       []string{"user:email"},
	RedirectURL:  "http://localhost:8080/api/auth/github/callback",
}

var GoogleConfig = &oauth2.Config{
	ClientID: os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Endpoint: google.Endpoint,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	RedirectURL: "http://localhost:8080/api/auth/google/callback",
}