package handler

import (
	"errors"
	"fmt"
	"github.com/VoRaX00/shortener/internal/handler/mocks"
	"github.com/VoRaX00/shortener/internal/service/shortener"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMainPage(t *testing.T) {
	mockShortenerService := mocks.NewShortenerService(t)
	h := New(slog.Default(), mockShortenerService)

	testCases := []struct {
		method       string
		url          string
		expectedCode int
		expectedBody string
	}{
		{
			method:       http.MethodGet,
			url:          "EwHXdJfB",
			expectedCode: http.StatusTemporaryRedirect,
			expectedBody: "<a href=\"http://example.com\">Temporary Redirect</a>.\n\n",
		},
		{
			method:       http.MethodPost,
			expectedCode: http.StatusCreated,
			expectedBody: "http://example.com/EwHXdJfB",
		},
		{
			method:       http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: ``,
		},
		{
			method:       http.MethodDelete,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: ``,
		},
		{
			method:       http.MethodPatch,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: ``,
		},
		{
			method:       http.MethodConnect,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: ``,
		},
		{
			method:       http.MethodHead,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: ``,
		},
		{
			method:       http.MethodOptions,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: ``,
		},
		{
			method:       http.MethodTrace,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: ``,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.method, func(t *testing.T) {
			if testCase.method == http.MethodGet {
				mockShortenerService.On("GetLink", mock.AnythingOfType("string")).Return("http://example.com", nil).Once()
			} else if testCase.method == http.MethodPost {
				mockShortenerService.On("Shorten", mock.AnythingOfType("string")).Return("EwHXdJfB", nil).Once()
			}

			r := httptest.NewRequest(testCase.method, fmt.Sprintf("/%s", testCase.url), nil)
			w := httptest.NewRecorder()
			h.mainPage(w, r)

			assert.Equal(t, testCase.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, testCase.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
		})
	}
}

func TestGetLink(t *testing.T) {
	mockShortenerService := mocks.NewShortenerService(t)
	h := New(slog.Default(), mockShortenerService)

	testCases := []struct {
		name            string
		id              string
		expectedCode    int
		expectedBody    string
		mockReturnError error
		mockReturnLink  string
	}{
		{
			name:            "Успешное получение ссылки по id",
			id:              "EwHXdJfB",
			mockReturnLink:  "http://example.com",
			mockReturnError: nil,
			expectedCode:    http.StatusTemporaryRedirect,
			expectedBody:    "<a href=\"http://example.com\">Temporary Redirect</a>.\n\n",
		},
		{
			name:            "Пустой id",
			id:              "",
			mockReturnLink:  "http://example.com",
			mockReturnError: nil,
			expectedCode:    http.StatusBadRequest,
		},
		{
			name:            "Не найдена ссылка по id",
			id:              "EwHXdJfB",
			mockReturnLink:  "",
			mockReturnError: shortener.ErrNotFound,
			expectedCode:    http.StatusNotFound,
		},
		{
			name:            "Ошибка сервиса",
			id:              "EwHXdJfB",
			mockReturnLink:  "",
			mockReturnError: errors.New(""),
			expectedCode:    http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.id != "" {
				mockShortenerService.On("GetLink", testCase.id).Return(testCase.mockReturnLink, testCase.mockReturnError).Once()
			}

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", testCase.id), nil)
			w := httptest.NewRecorder()
			h.getLink(w, r)

			assert.Equal(t, testCase.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, testCase.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
		})
	}
}

func TestShortener(t *testing.T) {
	mockShortenerService := mocks.NewShortenerService(t)
	h := New(slog.Default(), mockShortenerService)

	testCases := []struct {
		name            string
		link            string
		expectedCode    int
		expectedBody    string
		mockReturnError error
		mockReturnId    string
	}{
		{
			name:            "Успешное сокращение ссылки",
			link:            "http://yandex.com",
			expectedCode:    http.StatusCreated,
			expectedBody:    "http://example.com/EwHXdJfB",
			mockReturnError: nil,
			mockReturnId:    "EwHXdJfB",
		},
		{
			name:            "Ошибка в сервисе",
			link:            "http://yandex.com",
			expectedCode:    http.StatusInternalServerError,
			mockReturnError: errors.New("ошибка сервиса"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.link != "" {
				mockShortenerService.On("Shorten", testCase.link).Return(testCase.mockReturnId, testCase.mockReturnError).Once()
			}

			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testCase.link))
			w := httptest.NewRecorder()
			h.shortener(w, r)
			assert.Equal(t, testCase.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, testCase.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
		})
	}
}
