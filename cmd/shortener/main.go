package main

import (
	"github.com/Ippolid/shortLink/config"
	"github.com/Ippolid/shortLink/internal/app"
)

// функция main вызывается автоматически при запуске приложения
func main() {
	host, adr := config.ParseFlags()
	//fmt.Println(host, adr)
	k := app.NewDbase()
	s := app.New(&k, adr, host)
	if err := s.Start(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
