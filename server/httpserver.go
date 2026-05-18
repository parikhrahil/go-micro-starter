package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/parikhrahil/go-micro-starter/logger"

	"github.com/go-playground/validator/v10"
)

type s struct {
	server *http.Server
	notify chan error
	Logger logger.Logger
}

type HttpServerOpts struct {
	URL    string        `validate:"required"`
	Router http.Handler  `validate:"required"`
	Log    logger.Logger `validate:"required"`
}

func NewHTTPServer(opts *HttpServerOpts) (Server, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(opts); err != nil {
		return nil, err
	}

	router := opts.Router
	log := opts.Log

	server := &http.Server{
		Addr:    opts.URL,
		Handler: router,
	}
	return &s{
		server: server,
		notify: make(chan error, 1), // buffered to prevent leaking goroutines
		Logger: log,
	}, nil
}

func (s *s) serve() error {
	go func() {
		s.Logger.Info(fmt.Sprintf("Starting HTTP server on %s", s.server.Addr))
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.notify <- err
		}
		close(s.notify)
	}()
	return nil
}

func (s *s) shutdown(timeout time.Duration) error {
	s.Logger.Info("Shutting down HTTP Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return s.server.Close() // Force close if graceful shutdown fails
	}

	s.Logger.Info("HTTP server shut down gracefully")
	return nil
}

func (s *s) Run() {
	s.serve()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-s.notify:
		s.Logger.Fatal(fmt.Sprintf("Server start up failed: %v", err))

	case sig := <-interrupt:
		s.Logger.Info(fmt.Sprintf("Received signal: %v. Initiating teardown ...", sig))

		if err := s.shutdown(5 * time.Second); err != nil {
			s.Logger.Fatal(fmt.Sprintf("Server forced to shutdown with error: %v", err))
		}
	}
}
