package handlerserver

import (
	"github.com/Ippolid/shortLink/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/resty.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

//	func TestСreateLink(t *testing.T) {
//		db := NewDbase()
//		type want struct {
//			code        int
//			response    string
//			contentType string
//		}
//		tests := []struct {
//			name    string
//			request string
//			want    want
//		}{
//			{
//				name:    "positive test #1",
//				request: "https://ru.wikipedia.org/wiki/SHA-1",
//				want: want{
//					code:        201,
//					response:    "http://localhost:8080/b12a6809",
//					contentType: "text/plain",
//				},
//			},
//			{
//				name:    "positive test #2",
//				request: "https://github.com/Ippolid/shortLink/pulls?q=is%3Apr+is%3Aopen",
//				want: want{
//					code:        201,
//					response:    `http://localhost:8080/14603b1d`,
//					contentType: "text/plain",
//				},
//			},
//		}
//		for _, test := range tests {
//			t.Run(test.name, func(t *testing.T) {
//				server := New(&db) // создаём новый экземпляр сервера
//				handler := ValidationMiddleware(http.HandlerFunc(server.PostCreate))
//				request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.request))
//				request.Header.Set("Content-Type", "text/plain")
//				// создаём новый Recorder
//				w := httptest.NewRecorder()
//				handler.ServeHTTP(w, request)
//
//				res := w.Result()
//				// проверяем код ответа
//				assert.Equal(t, res.StatusCode, test.want.code)
//				// получаем и проверяем тело запроса
//				defer res.Body.Close()
//				resBody, err := io.ReadAll(res.Body)
//
//				require.NoError(t, err)
//				assert.Equal(t, string(resBody), test.want.response)
//				assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)
//			})
//		}
//	}
func TestCreateLink(t *testing.T) {
	db := storage.NewDbase()
	host := "localhost:8080"
	adr := "http://localhost:8080/"
	server := New(&db, adr, host, nil)
	r := server.newServer()
	ts := httptest.NewServer(r)

	defer ts.Close()

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

	client := resty.New()
	_, err := client.R().
		SetBody("https://test-url.com").
		Post(ts.URL + "/")
	require.NoError(t, err, "Ошибка при первом запросе `POST /`")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := client.R().
				SetBody(test.request).
				Post(ts.URL + "/")

			require.NoError(t, err)
			assert.Equal(t, test.want.code, resp.StatusCode())
			assert.Equal(t, test.want.response, resp.String())
			assert.Equal(t, test.want.contentType, resp.Header().Get("Content-Type"))
		})
	}
}

func TestGetLink(t *testing.T) {
	db := storage.NewDbase()
	host := "localhost:8080"
	adr := "http://localhost:8080/"
	server := New(&db, adr, host, nil)

	server.database.SaveLink([]byte("https://ru.wikipedia.org/wiki/SHA-1"), "b12a6809")
	server.database.SaveLink([]byte("https://github.com/Ippolid/shortLink/pulls?q=is%3Apr+is%3Aopen"), "14603b1d")

	r := server.newServer()
	ts := httptest.NewServer(r)
	defer ts.Close()

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
			request: "/b12a6809",
			want: want{
				code:     307,
				location: "https://ru.wikipedia.org/wiki/SHA-1",
			},
		},
		{
			name:    "positive test #2",
			request: "/14603b1d",
			want: want{
				code:     307,
				location: "https://github.com/Ippolid/shortLink/pulls?q=is%3Apr+is%3Aopen",
			},
		},
		{
			name:    "negative test #3",
			request: "/14603b1d123",
			want: want{
				code:     400,
				location: "",
			},
		},
	}

	client := resty.New().
		SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}))
	_, err := client.R().
		SetBody("https://ru.wikipedia.org/wiki/SHA-1").
		Post(ts.URL + "/")
	require.NoError(t, err, "Ошибка при первом запросе `POST /`")

	_, err = client.R().
		SetBody("https://github.com/Ippolid/shortLink/pulls?q=is%3Apr+is%3Aopen").
		Post(ts.URL + "/")
	require.NoError(t, err, "Ошибка при первом запросе `POST /`")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := client.R().
				Get(ts.URL + test.request)

			require.NoError(t, err)
			assert.Equal(t, test.want.code, resp.StatusCode())
			assert.Equal(t, test.want.location, resp.Header().Get("Location"))
		})
	}
}

func TestCreateLinkApi(t *testing.T) {
	db := storage.NewDbase()
	host := "localhost:8080"
	adr := "http://localhost:8080/"
	server := New(&db, adr, host, nil)
	r := server.newServer()
	ts := httptest.NewServer(r)
	defer ts.Close()

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
			request: `{"url":"https://ru.wikipedia.org/wiki/SHA-1"}`,
			want: want{
				code:        201,
				response:    `{"result":"http://localhost:8080/b12a6809"}`,
				contentType: "application/json",
			},
		},
		{
			name:    "positive test #2",
			request: `{"url":"https://github.com/Ippolid/shortLink/pulls?q=is%3Apr+is%3Aopen"}`,
			want: want{
				code:        201,
				response:    `{"result":"http://localhost:8080/14603b1d"}`,
				contentType: "application/json",
			},
		},
	}

	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"url":"https://test.com"}`).
		Post(ts.URL + "/api/shorten")
	require.NoError(t, err, "Ошибка при первом запросе `POST /`")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(test.request).
				Post(ts.URL + "/api/shorten")

			require.NoError(t, err)
			assert.Equal(t, test.want.code, resp.StatusCode())
			assert.JSONEq(t, test.want.response, resp.String())
			assert.Equal(t, test.want.contentType, resp.Header().Get("Content-Type"))
		})
	}
}
