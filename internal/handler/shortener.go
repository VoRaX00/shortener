package handler

import "net/http"

type ShortenerService interface {
	Shorten(url string) (string, error)
	GetLink(id string) (string, error)
}

func (h *Handler) mainPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.shortener(w, r)
	case http.MethodPost:
		h.rootPage(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) rootPage(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) shortener(w http.ResponseWriter, r *http.Request) {}
