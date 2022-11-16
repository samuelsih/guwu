package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samuelsih/guwu/service/auth"
)

type AuthDeps struct {
}

func AuthRoutes(deps AuthDeps) *chi.Mux {
	mux := chi.NewMux()

	mux.Get("/auth", generateRedirect())
	mux.Get("/auth/{provider}/callback", callback())

	return mux
}

func generateRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := r.URL.Query().Get("provider")
		encoder := json.NewEncoder(w)
		
		if provider == "" {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(map[string]any{
				"code": http.StatusBadRequest,
				"message": "unknown provider",
			})
			return
		}

		out := auth.OAuthLogin(r.Context(), provider)
		// w.Write([]byte(out.Link))
		http.Redirect(w, r, out.Link, http.StatusTemporaryRedirect)
	}
}

func callback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		
		provider := chi.URLParam(r, "provider")
		if provider == "" {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(map[string]any{
				"code": http.StatusBadRequest,
				"message": "unknown provider",
			})
			
			return
		}

		state := r.FormValue("state")
		if state == "" {
			w.WriteHeader(http.StatusUnauthorized)
			encoder.Encode(map[string]any{
				"code": http.StatusUnauthorized,
				"message": "unknown state",
			})

			return
		}

		code := r.FormValue("code")
		if code == "" {
			w.WriteHeader(http.StatusUnauthorized)
			encoder.Encode(map[string]any{
				"code": http.StatusUnauthorized,
				"message": "unknown code",
			})

			return
		}

		out := auth.Authorize(r.Context(), provider, state, code)
		encoder.Encode(out)
	}
}