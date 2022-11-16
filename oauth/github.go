package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/samuelsih/guwu/model"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var GithubConfig = &oauth2.Config{
	ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	Endpoint:     github.Endpoint,
	Scopes:       []string{"read:user", "user:gmail"},
	RedirectURL:  "http://localhost:8080/api/auth/github/callback",
}


type GithubProvider struct {
	state string
	config *oauth2.Config
}

func NewGithubProvider(state string) *GithubProvider {
	return &GithubProvider{
		state: state,
		config: &oauth2.Config{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			Endpoint:     github.Endpoint,
			Scopes:       []string{"read:user", "user:gmail"},
			RedirectURL:  "http://localhost:8080/api/auth/github/callback",
		},
	}
}

func (g *GithubProvider) URLRedirect() string {
	return g.config.AuthCodeURL(g.state)
}

type GithubUser struct {

}

func (g *GithubProvider) Authenticate(r *http.Request) (model.User, error) {
	var user model.User
	
	state := r.FormValue("state")
	if state != g.state {
		return user, fmt.Errorf("unknown state")
	}

	code := r.FormValue("code")
	if code == "" {
		return user, fmt.Errorf("unknown code")
	}

	token, err := g.config.Exchange(r.Context(), code)
	if err != nil {
		return user, err
	}

	if token == nil {
		return user, fmt.Errorf("token is nil")
	}

	client := g.config.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))

	if err != nil {
		return user, err
	}

	if resp == nil {
		return user, fmt.Errorf("response is nil")
	}

	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return user, err
	}

	if err := json.Unmarshal(responseBytes, &user); err != nil {
		return user, err
	}

	return user, nil
}