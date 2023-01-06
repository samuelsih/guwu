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

func setSessionCookie(w http.ResponseWriter, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name: "sid",
		Value: value,
		MaxAge: maxAge,
		SameSite: http.SameSiteLaxMode,
	})
}
