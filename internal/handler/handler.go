package handler

import "net/http"

type Handler struct {
	mux *http.ServeMux
}

func New() *Handler {
	return &Handler{
		mux: http.NewServeMux(),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) InitRoutes() {
	h.mux.HandleFunc("/", h.mainPage)
}
