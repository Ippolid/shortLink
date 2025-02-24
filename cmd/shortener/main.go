package main

import (
	"fmt"
	"github.com/Ippolid/shortLink/config"
	"github.com/Ippolid/shortLink/internal/app"
	"github.com/Ippolid/shortLink/internal/app/handlerserver"
	"github.com/Ippolid/shortLink/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	host, adr, path := config.ParseFlags()
	//fmt.Println(host, adr)

	if envRunHost := os.Getenv("SERVER_ADDRESS"); envRunHost != "" {
		host = envRunHost
	}

	if envRunAdr := os.Getenv("BASE_URL"); envRunAdr != "" {
		adr = envRunAdr
		adr += "/"
	}

	if envRunPath := os.Getenv("FILE_STORAGE_PATH"); envRunPath != "" {
		path = envRunPath
	}

	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}

	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}
	k := app.NewDbase()

	err := k.ReadLocal(path)
	if err != nil {
		panic(err)
	}

	s := handlerserver.New(&k, adr, host)

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := s.Start(); err != nil {
			panic(err)
		}
	}()

	// Обрабатываем сигналы завершения
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	fmt.Println("Shutting down server...")

	if err := k.WriteLocal(path); err != nil {
		fmt.Printf("Error writing local data: %v\n", err)
	} else {
		fmt.Println("Successfully wrote local data")
	}

}
