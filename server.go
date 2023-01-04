package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

const (
	idleTimeout  = 1 * time.Minute
	readTimeout  = 10 * time.Second
	writeTimeout = 30 * time.Second
	ctxTimeout   = 5 * time.Second
)

type Server struct {
	Router       *chi.Mux
	Dependencies BusinessDeps

	eg *errgroup.Group
}

type BusinessDeps struct {
	DB *sqlx.DB
}

func (s *Server) Run(addr string) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.Router,
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	shutdownError := make(chan error)

	go s.notifyShutdown(srv, shutdownError)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return <-shutdownError
}

func (s *Server) Shutdown(ctx context.Context) error {
	var egCtx context.Context
	s.eg, egCtx = errgroup.WithContext(ctx)

	s.eg.Go(func() error {
		return s.Dependencies.DB.Close()
	})

	select {
	case <-egCtx.Done():
		return egCtx.Err()

	default:
		return s.eg.Wait()
	}
}

func (s *Server) notifyShutdown(server *http.Server, shutdownErr chan<- error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	shutdownErr <- server.Shutdown(ctx)
}
