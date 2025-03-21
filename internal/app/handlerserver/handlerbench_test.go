package handlerserver

import (
	"github.com/Ippolid/shortLink/internal/app/storage"
	"github.com/stretchr/testify/require"
	"gopkg.in/resty.v1"
	"net/http/httptest"
	"strconv"
	"testing"
)

func BenchmarkServer_PostCreate(b *testing.B) {
	db := storage.NewDbase()
	host := "localhost:8080"
	adr := "http://localhost:8080/"
	server := New(&db, adr, host, nil)
	r := server.newServer()
	ts := httptest.NewServer(r)

	defer ts.Close()
	count := 40
	links := make([]string, count)
	for i := 0; i < count; i++ {
		links[i] = "https://ru.wikipedia.org/wiki/" + strconv.Itoa(i)
	}
	client := resty.New()
	_, err := client.R().
		SetBody("https://test-url.com").
		Post(ts.URL + "/")
	require.NoError(b, err, "Ошибка при первом запросе `POST /`")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		link := links[i%count] // крутимся по ссылкам
		resp, err := client.R().
			SetBody(link).
			Post(ts.URL + "/")
		if err != nil {
			b.Fatalf("Ошибка запроса: %v", err)
		}
		if resp.StatusCode() != 201 && resp.StatusCode() != 409 {
			b.Fatalf("неожиданный статус код: %d", resp.StatusCode())
		}
	}
}

//	func BenchmarkServer_GET(b *testing.B) {
//		db := storage.NewDbase()
//		host := "localhost:8080"
//		adr := "http://localhost:8080/"
//		server := New(&db, adr, host, nil)
//		r := server.newServer()
//		ts := httptest.NewServer(r)
//
//		defer ts.Close()
//		count := 40
//		links := make([]string, count)
//		for i := 0; i < count; i++ {
//			links[i] = "https://ru.wikipedia.org/wiki/" + strconv.Itoa(i)
//		}
//		client := resty.New()
//		_, err := client.R().
//			SetBody("https://test-url.com").
//			Post(ts.URL + "/")
//		require.NoError(b, err, "Ошибка при первом запросе `POST /`")
//
//		linksnew := make([]string, count)
//		for i := 0; i < count; i++ {
//			link := links[i] // крутимся по ссылкам
//			resp, err := client.R().
//				SetBody(link).
//				Post(ts.URL + "/")
//			if err != nil {
//				b.Fatalf("Ошибка запроса: %v", err)
//			}
//			if resp.StatusCode() != 201 && resp.StatusCode() != 409 {
//				b.Fatalf("неожиданный статус код: %d", resp.StatusCode())
//			}
//			linksnew[i] = resp.String()
//		}
//		b.ResetTimer()
//		for i := 0; i < b.N; i++ { // крутимся по ссылкам
//			resp, err := client.R().
//				Get(linksnew[i%count])
//			if err != nil {
//				b.Fatalf("Ошибка запроса: %v", err)
//			}
//			if resp.StatusCode() != 200 && resp.StatusCode() != 302 {
//				b.Fatalf("неожиданный статус код: %d", resp.StatusCode())
//			}
//		}
//	}
func BenchmarkServer_GET(b *testing.B) {
	db := storage.NewDbase()
	server := New(&db, "http://localhost:8080", "localhost:8080", nil)
	r := server.newServer()
	ts := httptest.NewServer(r) // создаём тестовый сервер
	defer ts.Close()

	count := 40
	links := make([]string, count)
	for i := 0; i < count; i++ {
		links[i] = "https://ru.wikipedia.org/wiki/" + strconv.Itoa(i)
	}

	client := resty.New()

	// Прогрев: создаём короткие ссылки через POST
	linksnew := make([]string, count)
	for i := 0; i < count; i++ {
		link := links[i]
		resp, err := client.R().
			SetBody(link).
			Post(ts.URL + "/")
		if err != nil {
			b.Fatalf("Ошибка запроса: %v", err)
		}
		if resp.StatusCode() != 201 && resp.StatusCode() != 409 {
			b.Fatalf("неожиданный статус код: %d", resp.StatusCode())
		}
		linksnew[i] = ts.URL + "/" + resp.String() // фиксируем короткие ссылки на ts.URL
	}

	b.ResetTimer()

	// Бенчмаркинг GET запросов
	for i := 0; i < b.N; i++ {
		resp, err := client.R().
			Get(linksnew[i%count])
		if err != nil {
			b.Fatalf("Ошибка запроса: %v", err)
		}
		if resp.StatusCode() != 200 && resp.StatusCode() != 400 {
			b.Fatalf("неожиданный статус код: %d", resp.StatusCode())
		}
	}
}
