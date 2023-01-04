package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/samuelsih/guwu/business/auth"
	"github.com/samuelsih/guwu/business/health"
	"github.com/samuelsih/guwu/pkg/response"
	"github.com/samuelsih/guwu/presentation"
)

func (s *Server) MountMiddleware() {
	s.Router.Use(middleware.Recoverer)
}

func (s *Server) MountHandlers() {
	s.authHandlers()
	s.healthCheckHandlers()

	s.notFound()
	s.methodNotAllowed()
}

func (s *Server) authHandlers() {
	authentication := auth.Deps{
		DB: s.Dependencies.DB,
	}

	s.Router.Post("/login", presentation.Post(authentication.Login))
	s.Router.Post("/register", presentation.Post(authentication.Register))
}

func (s *Server) healthCheckHandlers() {
	healthCheck := health.Deps{
		DB: s.Dependencies.DB,
	}

	s.Router.Get("/health", presentation.Get(healthCheck.Check, presentation.Config{}))
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
