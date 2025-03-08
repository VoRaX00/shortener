package main

import (
	"context"
	"github.com/VoRaX00/shortener/internal/app"
	"github.com/VoRaX00/shortener/internal/config"
	"github.com/VoRaX00/shortener/internal/handler"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

const serverConfig = "./config/server.yml"

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	application := setupApp(serverConfig)

	go application.MustStart()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	application.MustStop(context.Background())
}

func setupApp(configPath string) *app.App {
	cfg := config.MustConfig[app.Config](configPath)
	return app.NewApp(cfg.Addr, handler.New())
}
