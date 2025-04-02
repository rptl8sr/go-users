package main

import (
	"fmt"
	"log/slog"
	"os"

	"go-users/internal/app"
	"go-users/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	var logger *slog.Logger
	switch cfg.Log.Format {
	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: cfg.Log.Level,
		}))
	case "text":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: cfg.Log.Level,
		}))
	}

	a, err := app.New(cfg, logger)
	if err != nil {
		fmt.Printf("failed to create application: %v\n", err)
		os.Exit(1)
	}

	if err = a.Run(); err != nil {
		fmt.Printf("application error: %v\n", err)
		os.Exit(1)
	}
}
