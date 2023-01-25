package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/business/auth"
	"github.com/samuelsih/guwu/business/follow"
	"github.com/samuelsih/guwu/business/health"
	"github.com/samuelsih/guwu/pkg/redis"
	"github.com/samuelsih/guwu/pkg/response"
	pr "github.com/samuelsih/guwu/presentation"
)

func loadRoutes(r *chi.Mux, deps Dependencies) {
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	r.Use(middleware.Recoverer)

	redisClient := redis.NewClient(deps.Redis, "sessionid_")

	authRoutes(r, deps.DB, redisClient)
	followHandlers(r, deps.DB, redisClient)

	healthCheckHandlers(r, deps)
	notFound(r)
	methodNotAllowed(r)
}

func authRoutes(r *chi.Mux, db *sqlx.DB, rdb *redis.Client) {
	deps := auth.Deps{
		DB:             db,
		CreateSession:  rdb.SetJSON,
		DestroySession: rdb.Destroy,
	}

	r.Post("/register", pr.Post(deps.Register, pr.OnlyDecodeOpts))
	r.Post("/login", pr.Post(deps.Login, pr.SetSessionWithDecodeOpts))
	r.Delete("/logout", pr.Delete(deps.Logout, pr.GetterSetterSessionOpts))
}

func followHandlers(r *chi.Mux, db *sqlx.DB, rdb *redis.Client) {
	follow := follow.Deps{
		DB:             db,
		GetUserSession: rdb.GetJSON,
	}

	r.Post("/follow", pr.Post(follow.Follow, pr.GetSessionWithDecodeOpts))
}

func healthCheckHandlers(r *chi.Mux, deps Dependencies) {
	healthCheck := health.Deps{
		DB: deps.DB,
	}

	r.Get("/health", pr.Get(healthCheck.Check, pr.Opts{}))
}

func notFound(r *chi.Mux) {
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		res := map[string]any{
			"status_code": http.StatusNotFound,
			"msg":         http.StatusText(http.StatusNotFound),
		}

		if err := response.JSON(w, http.StatusNotFound, res); err != nil {
			log.Println(err)
		}
	})
}

func methodNotAllowed(r *chi.Mux) {
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		res := map[string]any{
			"status_code": http.StatusMethodNotAllowed,
			"msg":         http.StatusText(http.StatusMethodNotAllowed),
		}

		if err := response.JSON(w, http.StatusMethodNotAllowed, res); err != nil {
			log.Println(err)
		}
	})
}
