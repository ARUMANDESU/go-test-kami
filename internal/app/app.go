package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ARUMANDESU/go-test-kami/internal/api"
	"github.com/ARUMANDESU/go-test-kami/internal/config"
	"github.com/ARUMANDESU/go-test-kami/internal/service/reservation"
	"github.com/ARUMANDESU/go-test-kami/internal/storage/postgresql"
	"github.com/ARUMANDESU/go-test-kami/pkg/logger"
)

type App struct {
	log        *slog.Logger
	httpServer *http.Server
	storage    postgresql.Storage
}

func NewApp(ctx context.Context, cfg config.Config, logger *slog.Logger) *App {
	const op = "app.NewApp"
	log := logger.With("op", op)

	log.Info("creating new app")

	storage, err := postgresql.NewStorage(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	reservationService := reservation.NewService(logger, storage, storage)

	httpAPI := api.NewAPI(logger, reservationService)

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
	const op = "app.Stop"
	log := a.log.With("op", op)

	if err := a.httpServer.Shutdown(context.Background()); err != nil {
		log.Error("HTTP server shutdown error", logger.Err(err))
	}

	a.storage.Close()

	return nil
}
