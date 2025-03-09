package handler

import "net/http"

type Handler struct {
	mux     *http.ServeMux
	service ShortenerService
}

func New(service ShortenerService) *Handler {
	h := &Handler{
		mux:     http.NewServeMux(),
		service: service,
	}

	h.mux.HandleFunc("/", h.mainPage)
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
