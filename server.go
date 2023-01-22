package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rueian/rueidis"
	"github.com/samuelsih/guwu/pkg/logger"
)

const (
	idleTimeout  = 1 * time.Minute
	readTimeout  = 10 * time.Second
	writeTimeout = 30 * time.Second
	ctxTimeout   = 5 * time.Second
	shutdownTimeout = 30 * time.Second
)

type shutdownFunc func(ctx context.Context) error 

type Dependencies struct {
	DB *sqlx.DB
	Redis rueidis.Client
	// many more will come
}

func RunServer(router *chi.Mux, addr string, dependencies Dependencies) {
	server := http.Server {
		Addr:         addr,
		Handler:      router,
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}	

	loadRoutes(router, dependencies)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	done := make(chan struct{})
	
	listenOnShutdown(&server, quit, done, map[string]shutdownFunc {
		"shutdown server": func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},

		"shutdown db": func(_ context.Context) error {
			return dependencies.DB.Close()
		},

		"shutdown redis": func(_ context.Context) error {
			dependencies.Redis.Close()
			return nil
		},
	})

	logger.SysInfo("Serve on localhost" + addr)

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		logger.SysErr(err)
		return
	}

	<- done
	close(done)
	close(quit)
	logger.SysInfo("shutdown success")
}

func listenOnShutdown(server *http.Server, quit chan os.Signal, done chan struct{}, ops map[string]shutdownFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	var wg sync.WaitGroup
	
	go func ()  {
		<-quit

		logger.SysInfo("shutting down app")
	
		for opName, op := range ops {
			wg.Add(1)
			opName := opName
			op := op

			logger.SysInfo("on operation " + opName)

			go func ()  {
				defer wg.Done()

				if err := op(ctx); err != nil {
					logger.Errorf("error on %s: %v", opName, err)
					return
				}	

				logger.SysInfo("operation " + opName + " has success") 
			}()
		}

		done <- struct{}{}
	}()

	wg.Wait()
}