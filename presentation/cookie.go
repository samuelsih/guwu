package presentation

import (
	"fmt"
	"net/http"
)

func getSessionCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("sid")
	if err != nil {
		return "", fmt.Errorf("unknown session")
	}

	return cookie.Value, nil
}

func setSessionCookie(w http.ResponseWriter, cookieName, cookieValue string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    cookieValue,
		MaxAge:   maxAge,
		SameSite: http.SameSiteLaxMode,
	})
}
