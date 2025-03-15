package handler

import (
	"log/slog"
	"net/http"
)

type Handler struct {
	mux     *http.ServeMux
	log     *slog.Logger
	service ShortenerService
}

func New(log *slog.Logger, service ShortenerService) *Handler {
	h := &Handler{
		mux:     http.NewServeMux(),
		log:     log,
		service: service,
	}

	h.mux.HandleFunc("/", h.mainPage)
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
