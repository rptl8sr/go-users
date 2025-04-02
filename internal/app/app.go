package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-users/internal/api"
	"go-users/internal/config"
	"go-users/internal/database"
)

// The App represents the core application structure including configuration, logging, database, and the HTTP server.
type App struct {
	db     database.DB
	cfg    *config.Config
	logger *slog.Logger
	server *http.Server
}

// New initializes and returns a new App instance configured with the provided config and logger. Returns an error if setup fails.
func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	db, err := database.New(context.Background(), cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
		ReadTimeout:  time.Duration(cfg.HTTP.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.HTTP.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.HTTP.IdleTimeout) * time.Second,
	}

	return &App{
		cfg:    cfg,
		logger: logger,
		db:     db,
		server: server,
	}, nil
}

// Run starts the server, handles shutdown signals, and properly cleans up resources such as the database connection.
func (a *App) Run() error {
	serverErrors := make(chan error, 1)

	handler, err := api.NewHandler(a.cfg.OpenAPI, a.logger, a.db)
	if err != nil {
		return fmt.Errorf("failed to create handler: %w", err)
	}
	a.server.Handler = handler

	go func() {
		a.logger.Info("Starting server",
			"addr", a.server.Addr,
			"mode", a.cfg.App.Mode,
		)
		serverErrors <- a.server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		a.logger.Info("Shutting down server...")
	case errServe := <-serverErrors:
		if errServe != nil && !errors.Is(errServe, http.ErrServerClosed) {
			return fmt.Errorf("server encountered an error: %w", errServe)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = a.server.Shutdown(ctx); err != nil {
		a.logger.Error("Server forced to shutdown", "error", err)
	}

	a.db.Close()

	a.logger.Info("Server exiting")
	return nil
}
