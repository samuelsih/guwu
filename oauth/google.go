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
	"golang.org/x/oauth2/google"
)

// var googleConfig = &oauth2.Config{
// 	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
// 	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
// 	Endpoint:     google.Endpoint,
// 	Scopes: []string{
// 		"https://www.googleapis.com/auth/userinfo.email",
// 		"https://www.googleapis.com/auth/userinfo.profile",
// 	},
// 	RedirectURL: "http://localhost:8080/api/auth/google/callback",
// }

type GoogleProvider struct {
	state string
	config *oauth2.Config
}

func NewGoogleProvider(state string) *GoogleProvider {
	return &GoogleProvider{
		state: state,
		config: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Endpoint:     google.Endpoint,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			RedirectURL: "http://localhost:8080/api/auth/google/callback",
		},
	}
}

func (g *GoogleProvider) URLRedirect() string {
	return g.config.AuthCodeURL(g.state)
}

type GoogleUser struct {

}

func (g *GoogleProvider) Authenticate(r *http.Request) (model.User, error) {
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