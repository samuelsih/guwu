package guwu

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
	"github.com/samuelsih/guwu/config"
	"github.com/samuelsih/guwu/service"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	Router    *chi.Mux
	DB        *sqlx.DB
	SessionDB *redis.Client
}

func NewServer() *Server {
	s := &Server{
		Router:    chi.NewRouter(),
		DB:        config.ConnectAndInitCockroach(),
		SessionDB: config.NewRedis(),
	}
	return s
}

func (s *Server) load() {
	guest := service.Guest{DB: s.DB}
	user := service.User{DB: s.DB}

	s.Router.Get("/user/{username}", getWithParam(user.FindUser, "username"))
	s.Router.Post("/register", post(guest.Register))
	s.Router.Post("/login", post(guest.Login))

}

func (s *Server) Run(stop <-chan os.Signal) {
	s.Router.Use(
		s.Logger(),
		s.JSONResponse(),
	)

	s.load()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.Router,
	}

	go func() {
		log.Info().Msg("Listening on port :8080")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal().Msg("Cant start server: " + err.Error())
			}
		}
	}()

	<-stop

	log.Info().Msg("Shutting Down received")

	eg := &errgroup.Group{}

	eg.Go(func() error {
		err := s.DB.Close()

		if err != nil {
			log.Error().Msg("Error closing db: " + err.Error())
			return err
		}

		log.Info().Msg("Closing DB Success")
		return nil
	})

	if err := eg.Wait(); err != nil {
		log.Fatal().Msg("Error on closing database" + err.Error())
	}

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal().Msg("Error shutting down: " + err.Error())
	}

	log.Info().Msg("Shutdown complete")
}
