package main

import (
	"github.com/Ippolid/shortLink/internal/app"
)

// функция main вызывается автоматически при запуске приложения
func main() {
	k := app.NewDbase()
	s := app.New(&k)
	if err := s.Start(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
