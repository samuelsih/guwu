package guwu

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/bytedance/sonic/encoder"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"github.com/samuelsih/guwu/service"
)

func (Server) JSONResponse() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(f)
	}
}

func (Server) Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			wr := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()

			defer func() {
				t2 := time.Now()

				if rec := recover(); rec != nil {
					log.Error().
						Str("type", "error").
						Timestamp().
						Interface("recover_info", rec).
						Bytes("debug_stack", debug.Stack()).
						Msg("log system error")
					http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				log.Info().
					Str("type", "access").
					Timestamp().
					Fields(map[string]any{
						"remote_ip":  r.RemoteAddr,
						"url":        r.URL.Path,
						"proto":      r.Proto,
						"method":     r.Method,
						"user_agent": r.Header.Get("User-Agent"),
						"status":     wr.Status(),
						"latency_ms": float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
						"bytes_in":   r.Header.Get("Content-Length"),
						"bytes_out":  wr.BytesWritten(),
					}).
					Msg("incoming_request")
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func (Server) CookieExists(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("app_session")
		if err != nil {
			encoder.NewStreamEncoder(w).Encode(service.CommonResponse{
				StatusCode: http.StatusForbidden,
				Msg:        http.StatusText(http.StatusForbidden),
			})
			return
		}

		h(w, r)
	}
}

