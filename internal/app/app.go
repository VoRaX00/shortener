package app

import (
	"context"
	"errors"
	"net/http"
)

type Config struct {
	Addr string `yaml:"addr" env-required:"true"`
}

type App struct {
	server *http.Server
}

func NewApp(addr string, handler http.Handler) *App {
	return &App{
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (a *App) MustStart() {
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (a *App) MustStop(ctx context.Context) {
	if err := a.server.Shutdown(ctx); err != nil {
		panic(err)
	}
}
