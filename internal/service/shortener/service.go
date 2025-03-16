package shortener

import (
	"errors"
	"fmt"
	"github.com/VoRaX00/shortener/internal/storage"
	"log/slog"
)

var (
	ErrNotFound = errors.New("shortener: not found")
)

//go:generate mockery --name=Repository --output=./mocks --case=underscore
type Repository interface {
	Add(link string) (string, error)
	Get(id string) (string, error)
}

type Service struct {
	repository Repository
	log        *slog.Logger
}

func NewService(log *slog.Logger, repository Repository) *Service {
	return &Service{
		repository: repository,
		log:        log,
	}
}

func (s *Service) Shorten(url string) (string, error) {
	const op = "service.shortener.Shorten"
	log := s.log.With(
		slog.String("op", op),
		slog.String("link", url),
	)

	log.Info("adding link")
	id, err := s.repository.Add(url)
	if err != nil {
		log.Error("Error with adding link in db: ", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("successfully added link")
	return id, nil
}

func (s *Service) GetLink(id string) (string, error) {
	const op = "service.shortener.GetLink"
	log := s.log.With(
		slog.String("op", op),
		slog.String("id", id),
	)

	log.Info("getting link")
	link, err := s.repository.Get(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Error("error not found link")
			return "", fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		log.Error("error with getting link in db: ", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("successfully retrieved link")
	return link, nil
}
