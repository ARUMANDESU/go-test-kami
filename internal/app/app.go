package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ARUMANDESU/go-test-kami/internal/api"
	"github.com/ARUMANDESU/go-test-kami/internal/config"
	"github.com/ARUMANDESU/go-test-kami/pkg/logger"
)

type App struct {
	log        *slog.Logger
	httpServer *http.Server
}

func NewApp(cfg config.Config, logger *slog.Logger) *App {
	const op = "app.NewApp"
	log := logger.With("op", op)

	log.Info("creating new app")

	httpAPI := api.NewAPI(logger, nil) // TODO: add reservation service

	httpServer := httpAPI.HTTPServer(":" + cfg.HTTP.Port)

	return &App{
		log:        logger,
		httpServer: &httpServer,
	}
}

func (a App) Start() error {
	go func() {
		if err := a.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("HTTP server error", logger.Err(err))
		}
	}()
	return nil
}

func (a App) Stop() error {
	return a.httpServer.Shutdown(context.TODO())
}
