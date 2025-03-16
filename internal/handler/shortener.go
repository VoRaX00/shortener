package handler

import (
	"errors"
	"fmt"
	"github.com/VoRaX00/shortener/internal/service/shortener"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

//go:generate mockery --name=ShortenerService --output=./mocks --case=underscore
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

	h.log.Info("Getting link from shortener", "id", id)
	link, err := h.service.GetLink(id)
	if err != nil {
		if errors.Is(err, shortener.ErrNotFound) {
			h.log.Error("Unable to find link", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.log.Error("Unable to get link", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.log.Info("Successfully got link from shortener", "link", link)
	w.Header().Set("Location", link)
	http.Redirect(w, r, link, http.StatusTemporaryRedirect)
}

func (h *Handler) shortener(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
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

	shortLink := fmt.Sprintf("http://%s/%s", r.Host, id)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortLink))
	if err != nil {
		return
	}
}
