package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/samuelsih/guwu/business/auth"
	"github.com/samuelsih/guwu/business/follow"
	"github.com/samuelsih/guwu/business/health"
	"github.com/samuelsih/guwu/pkg/redis"
	"github.com/samuelsih/guwu/pkg/response"
	"github.com/samuelsih/guwu/presentation"
)

func (s *Server) MountMiddleware() {
	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))
	s.Router.Use(middleware.Recoverer)
}

func (s *Server) MountHandlers() {
	s.authHandlers()
	s.healthCheckHandlers()
	s.followHandlers()

	s.notFound()
	s.methodNotAllowed()
}

func (s *Server) authHandlers() {
	rdb := redis.NewClient(s.Dependencies.RedisDB, "sessionId_")

	authentication := auth.Deps{
		DB:             s.Dependencies.DB,
		CreateSession:  rdb.SetJSON,
		DestroySession: rdb.Destroy,
	}

	s.Router.Post("/register", presentation.Post(authentication.Register, presentation.OnlyDecodeOpts))
	s.Router.Post("/login", presentation.Post(authentication.Login, presentation.SetSessionWithDecodeOpts))
	s.Router.Delete("/logout", presentation.Delete(authentication.Logout, presentation.GetterSetterSessionOpts))
}

func (s *Server) followHandlers() {
	rdb := redis.NewClient(s.Dependencies.RedisDB, "sessionId_")

	follow := follow.Deps{
		DB:             s.Dependencies.DB,
		GetUserSession: rdb.GetJSON,
	}

	s.Router.Post("/follow", presentation.Post(follow.Follow, presentation.GetSessionWithDecodeOpts))
}

func (s *Server) healthCheckHandlers() {
	healthCheck := health.Deps{
		DB: s.Dependencies.DB,
	}

	s.Router.Get("/health", presentation.Get(healthCheck.Check, presentation.Opts{}))
}

func (s *Server) notFound() {
	s.Router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		res := map[string]any{
			"status_code": http.StatusNotFound,
			"msg":         http.StatusText(http.StatusNotFound),
		}

		if err := response.JSON(w, http.StatusNotFound, res); err != nil {
			log.Println(err)
		}
	})
}

func (s *Server) methodNotAllowed() {
	s.Router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		res := map[string]any{
			"status_code": http.StatusMethodNotAllowed,
			"msg":         http.StatusText(http.StatusMethodNotAllowed),
		}

		if err := response.JSON(w, http.StatusMethodNotAllowed, res); err != nil {
			log.Println(err)
		}
	})
}
