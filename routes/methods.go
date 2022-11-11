package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/samuelsih/guwu/config"

	"github.com/samuelsih/guwu/service"
)

type (
	PostHandler [T service.CommonInput, U service.CommonOutput] func(context.Context, *T) U
	GetHandler [U service.CommonOutput] func(context.Context) U
	GetWithParamHandler [U service.CommonOutput] func(context.Context, string) U
) 

func Get[out service.CommonOutput](path string, handler GetHandler[out]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output := handler(r.Context())
		
		encoder := json.NewEncoder(w)

		if output.CommonRes().StatusCode == 0 {
			w.WriteHeader(http.StatusOK)
			encoder.Encode(output)
			return
		}

		w.WriteHeader(output.CommonRes().StatusCode)
		encoder.Encode(output)
	}
}

func GetWithParam[out service.CommonOutput](path, param string, handler GetWithParamHandler[out]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlParam := chi.URLParam(r, param)
		output := handler(r.Context(), urlParam)
		encoder := json.NewEncoder(w)

		if output.CommonRes().StatusCode == 0 {
			w.WriteHeader(http.StatusOK)
			encoder.Encode(output)
			return
		}

		w.WriteHeader(output.CommonRes().StatusCode)
		encoder.Encode(output)
	}
}

func Post[T service.CommonInput, U service.CommonOutput](path string, handler PostHandler[T, U]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userSession, err := deps.readCookie(r)
		encoder := json.NewEncoder(w)
		decoder := json.NewDecoder(r.Body)

		if err != nil {
			encoder.Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        err.Error(),
			})
			return
		}
		
		var in T

		in.CommonReq().UserSession = userSession

		defer r.Body.Close()
		decoder.DisallowUnknownFields()

		err = decoder.Decode(&in)
		if err != nil {
			encoder.Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        err.Error(),
			})

			return
		}

		out := handler(r.Context(), &in)

		if out.CommonRes().SessionID != "" {
			if errCookie := deps.setCookie(w, out.CommonRes().SessionID); errCookie != nil {
				w.WriteHeader(http.StatusInternalServerError)
				encoder.Encode(out)
				return
			}
		}

		w.WriteHeader(out.CommonRes().StatusCode)
		encoder.Encode(out)
	}
}

type googleUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
}

func googleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := config.GoogleConfig.AuthCodeURL("thefuck?")
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func googleCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")
		if state != "thefuck?" {
			fmt.Fprintln(w, "unknown state")
			return
		}

		code := r.FormValue("code")
		if code == "" {
			fmt.Fprintln(w, "unknown code")
			return
		}

		token, err := config.GoogleConfig.Exchange(r.Context(), code)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		if token == nil {
			fmt.Fprintln(w, "token is nil")
			return
		}

		client := config.GoogleConfig.Client(r.Context(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		if resp == nil {
			fmt.Fprintln(w, "response is nil")
			return
		}

		defer resp.Body.Close()

		responseBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		var u googleUser
		if err := json.Unmarshal(responseBytes, &u); err != nil {
			fmt.Fprintln(w, "error in marshaling: " + err.Error())
			return
		}

		json.NewEncoder(w).Encode(u)
	}
}