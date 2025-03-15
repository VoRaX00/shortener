package main

import (
	"context"
	"fmt"
	"github.com/VoRaX00/shortener/internal/app"
	"github.com/VoRaX00/shortener/internal/config"
	"github.com/VoRaX00/shortener/internal/handler"
	"github.com/VoRaX00/shortener/internal/service/shortener"
	"github.com/VoRaX00/shortener/internal/storage/postgres"
	shortenerrepo "github.com/VoRaX00/shortener/internal/storage/postgres/shortener"
	_ "github.com/VoRaX00/shortener/migrations"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

const (
	postgresConfig = "./config/postgres_test.yml"
	serverConfig   = "./config/server.yml"
	loggerConfig   = "./config/logger.yml"
	migrations     = "./migrations"
)

func main() {
	logger := setupLogger(loggerConfig)
	repository := setupPostgres(postgresConfig)
	defer func() {
		err := repository.Stop()
		if err != nil {
			fmt.Println(err)
		}
	}()

	shortenerService := shortener.NewService(logger, repository)
	application := setupApp(serverConfig, logger, shortenerService)

	logger.Info("starting shortener service")
	go application.MustStart()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("shortener service shutting down")
	application.MustStop(context.Background())
	logger.Info("shortener service stopped")
}

func setupApp(configPath string, log *slog.Logger, service handler.ShortenerService) *app.App {
	cfg := config.MustConfig[app.Config](configPath)
	h := handler.New(log, service)
	return app.NewApp(cfg.Addr, h)
}

func setupPostgres(configPath string) *shortenerrepo.Repository {
	cfg := config.MustConfig[postgres.Config](configPath)

	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode))
	if err != nil {
		panic(err)
	}

	if err = goose.Up(db.DB, migrations); err != nil {
		panic(err)
	}

	repo := shortenerrepo.NewRepository(db)
	return repo
}

func setupLogger(configPath string) *slog.Logger {
	cfg := config.MustConfig[config.Logger](configPath)

	var logger *slog.Logger

	switch cfg.Env {
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}
	return logger
}
