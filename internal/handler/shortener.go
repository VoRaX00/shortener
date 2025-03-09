package handler

import (
	"errors"
	"github.com/VoRaX00/shortener/internal/service/shortener"
	"io/ioutil"
	"net/http"
	"strings"
)

type ShortenerService interface {
	Shorten(url string) (string, error)
	GetLink(id string) (string, error)
}

func (h *Handler) mainPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getLink(w, r)
	case http.MethodPost:
		h.shortener(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getLink(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link, err := h.service.GetLink(id)
	if err != nil {
		if errors.Is(err, shortener.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", link)
	http.Redirect(w, r, link, http.StatusTemporaryRedirect)
}

func (h *Handler) shortener(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link := string(bodyBytes)
	id, err := h.service.Shorten(link)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(id))
	if err != nil {
		return
	}
}
