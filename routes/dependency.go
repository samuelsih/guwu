package routes

import (
	"errors"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/samuelsih/guwu/model"
)

type dependencies struct {
	Session model.SessionDeps
	GoogleState string
	GithubState string

	securer securecookie.SecureCookie
}

var deps dependencies

func InitDependency(sessionDeps model.SessionDeps, googleState, githubState string) {
	deps = dependencies{
		Session: sessionDeps,
		GoogleState: googleState,
		GithubState: githubState,
	}
}

func(d *dependencies) readCookie(r *http.Request) (model.Session, error) {
	cookie, err := r.Cookie("app_session")
	if err != nil || cookie.Value == "" {
		return model.Session{}, errors.New(http.StatusText(http.StatusBadRequest))
	}

	var sessionID string

	err = d.securer.Decode("app_session", cookie.Value, &sessionID)
	if err != nil {
		return model.Session{}, errors.New(http.StatusText(http.StatusBadRequest))
	}

	user, err := d.Session.Get(r.Context(), sessionID)

	if err != nil {
		return model.Session{}, errors.New(http.StatusText(http.StatusBadRequest))
	}

	return user, nil
}

func(d *dependencies) setCookie(w http.ResponseWriter, value string) error {
	encoded, err := d.securer.Encode("app_session", value)
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     "app_session",
		Value:    encoded,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		MaxAge:   24 * 3600,
	}

	http.SetCookie(w, &cookie)
	return nil
}