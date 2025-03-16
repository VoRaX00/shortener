package shortener

import (
	"errors"
	"github.com/VoRaX00/shortener/internal/service/shortener/mocks"
	"github.com/VoRaX00/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestShorten(t *testing.T) {
	mockRepository := mocks.NewRepository(t)
	service := NewService(slog.Default(), mockRepository)

	testErr := errors.New("error with repo")
	testCases := []struct {
		name          string
		url           string
		mockReturnId  string
		mockReturnErr error
		expectedId    string
		expectedErr   error
	}{
		{
			name:          "Успешное сокращение ссылки",
			url:           "http://example.com",
			mockReturnId:  "EwHXdJfB",
			mockReturnErr: nil,
			expectedId:    "EwHXdJfB",
			expectedErr:   nil,
		},
		{
			name:          "Ошибка в репозитории",
			url:           "http://example.com",
			mockReturnId:  "",
			mockReturnErr: testErr,
			expectedId:    "",
			expectedErr:   testErr,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.On("Add", tt.url).Return(tt.mockReturnId, tt.mockReturnErr).Once()

			id, err := service.Shorten(tt.url)
			assert.Equal(t, tt.expectedId, id)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestGetLink(t *testing.T) {
	mockRepository := mocks.NewRepository(t)
	service := NewService(slog.Default(), mockRepository)

	testErr := errors.New("error with repo")
	testCases := []struct {
		name           string
		id             string
		mockReturnLink string
		mockReturnErr  error
		expectedLink   string
		expectedErr    error
	}{
		{
			name:           "Успешное получение ссылки",
			id:             "EwHXdJfB",
			mockReturnLink: "http://example.com",
			mockReturnErr:  nil,
			expectedLink:   "http://example.com",
			expectedErr:    nil,
		},
		{
			name:           "Ссылка не найдена",
			id:             "EwHXdJfB",
			mockReturnLink: "",
			mockReturnErr:  storage.ErrNotFound,
			expectedLink:   "",
			expectedErr:    ErrNotFound,
		},
		{
			name:           "Ошибка в репозитории",
			id:             "EwHXdJfB",
			mockReturnLink: "",
			mockReturnErr:  testErr,
			expectedLink:   "",
			expectedErr:    testErr,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.On("Get", tt.id).Return(tt.mockReturnLink, tt.mockReturnErr).Once()
			link, err := service.GetLink(tt.id)
			assert.Equal(t, tt.expectedLink, link)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
