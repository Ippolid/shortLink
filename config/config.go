// package config
//
// import (
//
//	"flag"
//	"net/url"
//
// )
//
//	func ParseFlags() (string, string) {
//		host := flag.String("a", "localhost:8080", "host")
//		adr := flag.String("b", "http://localhost:8080/", "adress")
//		flag.Parse()
//		return *host, *adr
//	}
package config

import (
	"flag"
	"net/url"
	"os"
)

func ParseFlags() (string, string, string) {
	host := flag.String("a", "localhost:8080", "host")
	baseURL := flag.String("b", "http://localhost:8080/", "base URL")
	path := flag.String("f", "./tmp/short-url-db.json", "base URL")
	flag.Parse()

	// Проверяем и нормализуем baseURL
	if u, err := url.Parse(*baseURL); err == nil {
		// Убеждаемся, что URL заканчивается на /
		if u.Path == "" {
			u.Path = "/"
		}
		return *host, u.String(), *path
	}
	if *path == "" {
		*path = "./tmp/short-url-db.json"
	}

	if _, err := os.Stat(*path); os.IsNotExist(err) {
		// Создаем файл, если он не существует
		file, err := os.Create(*path)
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	return *host, "http://localhost:8080/", *path
}
