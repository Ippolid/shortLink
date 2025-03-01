package config

import (
	"flag"
	"net/url"
	"os"
)

// ParseFlags обрабатывает флаги и переменные окружения
func ParseFlags() (string, string, string, string) {
	// Значения по умолчанию
	defaultHost := "localhost:8080"
	defaultBaseURL := "http://localhost:8080/"
	defaultPath := "/tmp/short-url-db.json"
	defaultDb := "postgres://postgres:1234@localhost:5432/shorty"

	// Флаги командной строки
	host := flag.String("a", defaultHost, "Адрес сервера")
	baseURL := flag.String("b", defaultBaseURL, "Базовый URL")
	path := flag.String("f", defaultPath, "Путь к файлу хранения URL")
	db := flag.String("d", defaultDb, "Путь к базе данных")

	flag.Parse()

	// Приоритет конфигурации: env -> flag -> default
	if envHost := os.Getenv("SERVER_ADDRESS"); envHost != "" {
		*host = envHost
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		*baseURL = envBaseURL + "/"
	}
	if envPath := os.Getenv("FILE_STORAGE_PATH"); envPath != "" {
		*path = envPath
	}

	if envDb := os.Getenv("DATABASE_DSN"); envDb != "" {
		*db = envDb
	}

	if u, err := url.Parse(*baseURL); err == nil {
		// Убеждаемся, что URL заканчивается на /
		if u.Path == "" {
			u.Path = "/"
		}
		return *host, u.String(), *path, *db
	}

	// Если путь пустой, отключаем запись на диск
	if *path == "" {
		return *host, *baseURL, "", *db
	}

	return *host, *baseURL, *path, *db
}
