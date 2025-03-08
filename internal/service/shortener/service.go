package shortener

import (
	"errors"
	"github.com/VoRaX00/shortener/internal/storage"
)

var (
	ErrNotFound = errors.New("shortener: not found")
)

type Repository interface {
	Add(link string) (string, error)
	Get(id string) (string, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Shorten(url string) (string, error) {
	id, err := s.repository.Add(url)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Service) GetLink(id string) (string, error) {
	link, err := s.repository.Get(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}
	return link, nil
}
