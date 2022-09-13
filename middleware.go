package guwu

import "net/http"

func JSONResponse() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {	
		f := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(f)
	}
}