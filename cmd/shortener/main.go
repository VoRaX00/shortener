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
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
	"syscall"
)

const postgresConfig = "./config/postgres.yml"
const serverConfig = "./config/server.yml"

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	repository := setupPostgres(postgresConfig)
	defer func() {
		err := repository.Stop()
		if err != nil {
			fmt.Println(err)
		}
	}()

	shortenerService := shortener.NewService(repository)
	application := setupApp(serverConfig, shortenerService)

	go application.MustStart()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	application.MustStop(context.Background())
}

func setupApp(configPath string, service handler.ShortenerService) *app.App {
	cfg := config.MustConfig[app.Config](configPath)
	h := handler.New(service)
	return app.NewApp(cfg.Addr, h)
}

func setupPostgres(configPath string) *shortenerrepo.Repository {
	cfg := config.MustConfig[postgres.Config](configPath)
	cfg.Password = os.Getenv("POSTGRES_PASSWORD")

	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode))
	if err != nil {
		panic(err)
	}

	repo := shortenerrepo.NewRepository(db)
	return repo
}
