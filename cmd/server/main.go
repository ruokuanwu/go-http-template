package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/go-http-template/internal/config"
	"example.com/go-http-template/internal/controller"
	"example.com/go-http-template/internal/db"
	"example.com/go-http-template/internal/middleware"
	logpkg "example.com/go-http-template/internal/pkg/logger"
	"example.com/go-http-template/internal/pkg/response"
	"example.com/go-http-template/internal/repository"
	"example.com/go-http-template/internal/router"
	"example.com/go-http-template/internal/service"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/samber/oops"
	"go.uber.org/zap"
)

// @title Go HTTP Template API
// @version 1.0
// @description 用户认证 API。成功响应 code 为 0；失败响应仅包含 code 和 message。
// @BasePath /api
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
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	database, err := db.Open(cfg.DatabasePath)
	if err != nil {
		return err
	}
	connection, err := database.DB()
	if err != nil {
		return oops.Wrapf(err, "get database connection")
	}
	defer connection.Close()

	authService, err := service.NewAuth(repository.NewUser(database))
	if err != nil {
		return err
	}
	store := sqlite3store.New(connection)
	defer store.StopCleanup()
	sessions := scs.New()
	sessions.Store = store
	sessions.Lifetime = cfg.SessionLifetime
	sessions.IdleTimeout = cfg.SessionIdleTimeout
	sessions.Cookie.Name = "session_id"
	sessions.Cookie.HttpOnly = true
	sessions.Cookie.Secure = cfg.CookieSecure
	sessions.Cookie.SameSite = http.SameSiteLaxMode
	sessions.Cookie.Persist = true
	sessions.ErrorFunc = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Error("session persistence failed", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(response.ErrorResponse{Code: response.CodeInternalError, Message: "服务器内部错误"}); err != nil {
			logger.Error("write session error response", zap.Error(err))
		}
	}
	protection := http.NewCrossOriginProtection()
	for _, origin := range cfg.TrustedOrigins {
		if err := protection.AddTrustedOrigin(origin); err != nil {
			return oops.Wrapf(err, "configure trusted origin")
		}
	}
	auth := controller.NewAuth(authService, sessions)
	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           middleware.Session(sessions, router.New(auth, protection, logger, cfg.WebDir, cfg.DEBUG)),
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
