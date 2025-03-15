package handler

import (
	"fmt"
	"github.com/VoRaX00/shortener/internal/handler/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
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
			url:          "link",
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

}

func TestShortener(t *testing.T) {

}
