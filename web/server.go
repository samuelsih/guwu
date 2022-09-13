package web

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	router *chi.Mux
}

func InitServer() *Server {
	return &Server{
		router: chi.NewRouter(),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}