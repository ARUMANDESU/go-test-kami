package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ARUMANDESU/go-test-kami/internal/app"
	"github.com/ARUMANDESU/go-test-kami/internal/config"
	"github.com/ARUMANDESU/go-test-kami/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file")
	}

	cfg := config.MustLoad()

	log, close := logger.Setup(cfg.Env)
	defer close()

	log.Info("starting the app", slog.Attr{Key: "env", Value: slog.StringValue(cfg.Env)})

	application := app.NewApp(context.Background(), *cfg, log)

	go application.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))
	application.Stop()
	log.Info("application stopped", slog.String("signal", sign.String()))
}
