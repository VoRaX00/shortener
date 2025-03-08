package handler

import "net/http"

type Handler struct {
	mux     *http.ServeMux
	service ShortenerService
}

func New(service ShortenerService) *Handler {
	return &Handler{
		mux:     http.NewServeMux(),
		service: service,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.HandleFunc("/", h.mainPage)
	h.mux.ServeHTTP(w, r)
}
