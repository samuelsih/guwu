package presentation

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func getSessionCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return "", nil

		default:
			log.Println(err)
			return "", fmt.Errorf("can't get your session")
		}
	}

	return cookie.Value, nil
}
