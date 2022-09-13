package oauth

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var GithubConfig = &oauth2.Config{
	ClientID: os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	Endpoint: github.Endpoint,
	Scopes: []string{"read:user", "user:gmail"},
	RedirectURL: "http://localhost:3000/api/v1/oauth2/github/callback",
}

var OAuthStateString = "random"