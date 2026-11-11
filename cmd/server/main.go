package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/go-http-template/internal/config"
	"example.com/go-http-template/internal/controller"
	"example.com/go-http-template/internal/db"
	logpkg "example.com/go-http-template/internal/pkg/logger"
	"example.com/go-http-template/internal/repository"
	"example.com/go-http-template/internal/router"
	"example.com/go-http-template/internal/service"
	"github.com/samber/oops"
	"go.uber.org/zap"
)

func main() {
	logger, err := logpkg.New()
	if err != nil {
		panic(err)
	}
	if err := run(logger); err != nil {
		logger.Error("server stopped", zap.Any("error", err))
		_ = logger.Sync()
		os.Exit(1)
	}
	_ = logger.Sync()
}

func run(logger *zap.Logger) error {
	cfg := config.Load()
	database, err := db.Open(cfg.DatabasePath)
	if err != nil {
		return err
	}
	connection, err := database.DB()
	if err != nil {
		return oops.Wrapf(err, "get database connection")
	}
	defer connection.Close()

	todos := controller.NewTodo(service.NewTodo(repository.NewTodo(database)))
	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           router.New(todos, logger, cfg.WebDir),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	failed := make(chan error, 1)
	go func() { failed <- server.ListenAndServe() }()
	logger.Info("server starting", zap.String("address", cfg.Address))
	select {
	case err := <-failed:
		if !errors.Is(err, http.ErrServerClosed) {
			return oops.Wrapf(err, "serve http")
		}
		return nil
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return oops.Wrapf(server.Shutdown(shutdown), "shutdown http server")
	}
}
