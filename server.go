package guwu

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
	"github.com/samuelsih/guwu/config"
	"github.com/samuelsih/guwu/mail"
	"github.com/samuelsih/guwu/model"
	"github.com/samuelsih/guwu/service"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	Router    *chi.Mux
	DB        *sqlx.DB
	SessionDB *redis.Client
	Mailer *mail.Mailer
}

func NewServer(production bool) *Server {
	if production {
		s := &Server{
			Router:    chi.NewRouter(),
			DB:        config.ConnectPostgres(os.Getenv("COCKROACH_DSN")),
			SessionDB: config.NewRedis(os.Getenv("REDIS_URL")),
			Mailer: mail.NewMailer(os.Getenv("MAILER_HOST"), os.Getenv("MAILER_ADDR"), os.Getenv("MAILER_EMAIL"), os.Getenv("MAILER_PASSWORD")),
		}

		config.MigrateAll(s.DB)

		return s
	}

	s := &Server{
		Router:    chi.NewRouter(),
		DB:        config.ConnectPostgres(""),
		SessionDB: config.NewRedis(""),
		Mailer: mail.NewMailer(os.Getenv("MAILER_HOST"), os.Getenv("MAILER_ADDR"), os.Getenv("MAILER_EMAIL"), os.Getenv("MAILER_PASSWORD")),
	}

	if err := config.MigrateAll(s.DB); err != nil {
		log.Fatal().Msg("Cant migrate: " + err.Error())
	}

	go s.Mailer.Listen()

	return s
}

func (s *Server) load() {
	session := model.SessionDeps{Conn: s.SessionDB}

	guest := service.Guest{DB: s.DB, SessionDB: s.SessionDB}
	posts := service.Post{DB: s.DB}

	s.Router.Post("/register", loginOrRegister(guest.Register))
	s.Router.Post("/login", loginOrRegister(guest.Login))
	s.Router.Post("/logout", logout(guest.Logout))

	s.Router.Get("/timeline", get(posts.Timeline))
	s.Router.Post("/post", post(session, posts.Insert))
	s.Router.Put("/post/{id}", put(session, "id", posts.Edit))
}

func (s *Server) Run(stop <-chan os.Signal) {
	s.Router.Use(
		middleware.RequestID,
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

	s.Mailer.Close()

	eg.Go(func() error {
		err := s.DB.Close()

		if err != nil {
			log.Error().Msg("Error closing db: " + err.Error())
			return err
		}

		log.Info().Msg("Closing DB Success")
		return nil
	})

	eg.Go(func() error {
		err := s.SessionDB.Close()

		if err != nil {
			log.Error().Msg("Error closing session db: " + err.Error())
			return err
		}

		log.Info().Msg("Closing Session DB Success")
		return nil
	})

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal().Msg("Error shutting down: " + err.Error())
	}

	if err := eg.Wait(); err != nil {
		log.Fatal().Msg("Error on closing database" + err.Error())
	}

	log.Info().Msg("Shutdown complete")
}
