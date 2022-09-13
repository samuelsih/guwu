package guwu

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/config"
	"github.com/samuelsih/guwu/service"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	Router *chi.Mux
	DB *sqlx.DB
}

func NewServer() *Server {
	s := &Server{
		Router: chi.NewRouter(),
		DB: config.ConnectAndInitCockroach(),
	}
	return s
}

func (s *Server) load() {
	guest := service.Guest{DB: s.DB}
	POST(s.Router, "/register", guest.Register)
	POST(s.Router, "/login", guest.Login)
}

func (s *Server) Run(stop <-chan os.Signal) {
	s.Router.Use(
		JSONResponse(),
		middleware.Recoverer,
	)

	s.load()

	srv := &http.Server{
		Addr: ":8080",
		Handler: s.Router,
	}

	go func ()  {
		log.Println("Listening on port :8080")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}	
	}()

	<- stop

	log.Println("Shutting Down received")

	eg := &errgroup.Group{}

	eg.Go(func() error {
		err := s.DB.Close()

		if err != nil {
			log.Println("Error closing db:", err)
			return err
		}

		log.Println("Closing DB Success")
		return nil
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}

	log.Println("Shutdown complete")
}