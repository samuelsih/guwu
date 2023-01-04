package presentation

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	b "github.com/samuelsih/guwu/business"
	"github.com/samuelsih/guwu/pkg/request"
	"github.com/samuelsih/guwu/pkg/response"
)

type handler[inType b.CommonInput, outType b.CommonOutput] func(ctx context.Context, in inType) outType

type Config struct {
	GetSessionCookie bool
	SetSessionCookie bool
}

func Get[inType b.CommonInput, outType b.CommonOutput](handle handler[inType, outType], config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in inType

		if config.GetSessionCookie {
			sessionID, err := getSessionCookie(r)
			if err != nil {
				encodeErr := response.JSON(w, 400, err)

				if encodeErr != nil {
					log.Printf("presentation.Get: %v", err)
				}

				return
			}

			in.CommonReq().SessionID = sessionID
		}

		out := handle(r.Context(), in)

		if err := response.JSON(w, out.CommonRes().StatusCode, &out); err != nil {
			log.Printf("presentation.Get: %v", err)
			return
		}
	}
}

func Post[inType any, outType b.CommonOutput](handler func(ctx context.Context, in inType) outType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in inType

		if err := request.Decode(w, r, &in); err != nil {
			encodeErr := response.JSON(w, 400, err)

			if encodeErr != nil {
				log.Printf("presentation.Post: %v", err)
			}

			return
		}

		defer r.Body.Close()

		out := handler(r.Context(), in)

		if err := response.JSON(w, out.CommonRes().StatusCode, out); err != nil {
			log.Printf("presentation.Post: %v", err)
			return
		}
	}
}

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
