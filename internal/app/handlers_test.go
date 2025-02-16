package app

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestСreateLink(t *testing.T) {
	db := NewDbase()
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "positive test #1",
			request: "https://ru.wikipedia.org/wiki/SHA-1",
			want: want{
				code:        201,
				response:    "http://localhost:8080/b12a6809",
				contentType: "text/plain",
			},
		},
		{
			name:    "positive test #2",
			request: "https://github.com/Ippolid/shortLink/pulls?q=is%3Apr+is%3Aopen",
			want: want{
				code:        201,
				response:    `http://localhost:8080/14603b1d`,
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := New(&db) // создаём новый экземпляр сервера
			handler := ValidationMiddleware(http.HandlerFunc(server.PostCreate))
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.request))
			request.Header.Set("Content-Type", "text/plain")
			// создаём новый Recorder
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, res.StatusCode, test.want.code)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, string(resBody), test.want.response)
			assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)
		})
	}
}
func TestGetLink(t *testing.T) {
	db := NewDbase()
	server := New(&db)
	server.database.SaveLink([]byte("https://ru.wikipedia.org/wiki/SHA-1"), "b12a6809")
	server.database.SaveLink([]byte("https://github.com/Ippolid/shortLink/pulls?q=is%3Apr+is%3Aopen"), "14603b1d")
	type want struct {
		code     int
		location string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "positive test #1",
			request: "http://localhost:8080/b12a6809",
			want: want{
				code:     307,
				location: "https://ru.wikipedia.org/wiki/SHA-1",
			},
		},
		{
			name:    "positive test #2",
			request: "http://localhost:8080/14603b1d",
			want: want{
				code:     307,
				location: "https://github.com/Ippolid/shortLink/pulls?q=is%3Apr+is%3Aopen",
			},
		},
		{
			name:    "negative test #3",
			request: "http://localhost:8080/14603b1d123",
			want: want{
				code:     400,
				location: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// создаём новый экземпляр сервера
			handler := ValidationMiddleware(http.HandlerFunc(server.GetID))
			request := httptest.NewRequest(http.MethodGet, test.request, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()
			// проверяем код ответа
			assert.Equal(t, res.StatusCode, test.want.code)
			// проверяем заголовок Location
			location := res.Header.Get("Location")
			assert.Equal(t, location, test.want.location)
		})
	}
}
