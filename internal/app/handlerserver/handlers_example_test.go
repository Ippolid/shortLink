package handlerserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Ippolid/shortLink/internal/app/storage"
	"github.com/Ippolid/shortLink/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

// Пример использования хендлера GetUserURLs для получения ссылок пользователя
func ExampleServer_GetUserURLs() {
	// Инициализация сервера для примера
	db := &storage.Dbase{}
	server := New(db, "http://localhost:8080/", "localhost:8080", nil)

	// Создаем тестовый Gin контекст
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Имитируем авторизованного пользователя
	c.Set("user_id", "test-user-id")

	// Вызываем хендлер
	server.GetUserURLs(c)

	// Выводим статус код и тело ответа
	fmt.Printf("Status Code: %d\n", w.Code)
	// Вывод в реальном применении будет зависеть от имеющихся данных
	// Output:
	// Status Code: 204
}

// Пример использования хендлера PostBatch для создания нескольких ссылок
func ExampleServer_PostBatch() {
	// Инициализация сервера для примера
	db := &storage.Dbase{}
	server := New(db, "http://localhost:8080/", "localhost:8080", nil)

	// Создаем тестовый Gin контекст
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Имитируем авторизованного пользователя
	c.Set("user_id", "test-user-id")

	// Подготавливаем входные данные
	batchReq := []models.PostBatchReq{
		{ID: "1", URL: "https://example.com"},
		{ID: "2", URL: "https://example.org"},
	}
	reqBody, _ := json.Marshal(batchReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer(reqBody))

	// Вызываем хендлер
	server.PostBatch(c)

	// Выводим статус код
	fmt.Printf("Status Code: %d\n", w.Code)
	// Вывод в реальном применении будет содержать созданные короткие ссылки
	// Output:
	// Status Code: 201
}

// Пример использования хендлера DeleteLinks для удаления ссылок пользователя
func ExampleServer_DeleteLinks() {
	// Инициализация сервера для примера
	db := &storage.Dbase{}
	// Для этого примера нам нужен mock базы данных, так как DeleteLinks использует Db
	mockDB := &storage.DataBase{}
	server := New(db, "http://localhost:8080/", "localhost:8080", mockDB)

	// Создаем тестовый Gin контекст
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Имитируем авторизованного пользователя
	c.Set("user_id", "test-user-id")

	// Подготавливаем входные данные - список идентификаторов для удаления
	linksToDelete := []string{"abc123", "def456"}
	reqBody, _ := json.Marshal(linksToDelete)
	c.Request = httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBuffer(reqBody))

	// Вызываем хендлер
	server.DeleteLinks(c)

	// Выводим статус код
	fmt.Printf("Status Code: %d\n", w.Code)
	// Output:
	// Status Code: 202
}
