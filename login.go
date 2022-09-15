package guwu

import (
	"net/http"

	"github.com/samuelsih/guwu/oauth"
)

func Login(w http.ResponseWriter, r *http.Request) {
	url := oauth.GithubConfig.AuthCodeURL(oauth.OAuthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func Callback(w http.ResponseWriter, r *http.Request) {

}
