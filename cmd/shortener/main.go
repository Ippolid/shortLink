package main

import (
	"context"
	"fmt"
	"github.com/Ippolid/shortLink/internal/app/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ippolid/shortLink/config"
	"github.com/Ippolid/shortLink/internal/app/handlerserver"
	"github.com/Ippolid/shortLink/internal/logger"
)

func main() {
	var server *handlerserver.Server
	var db = storage.NewDbase()
	// Получаем параметры конфигурации
	host, baseURL, storagePath, dboopen := config.ParseFlags()

	// Инициализация логгера
	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}
	if dboopen == "" {
		// Создаём базу данных
		db := storage.NewDbase()

		// Загружаем данные из файла, если файл указан
		if err := db.ReadLocal(storagePath); err != nil {
			log.Fatalf("Ошибка загрузки данных: %v", err)
		}
		server = handlerserver.New(&db, baseURL, host, nil)
	} else {

		// Создаём базу данных postresql
		db1, err := storage.Connect(dboopen)
		if err != nil {
			log.Fatal("open", err)
		}

		datab, err := storage.NewDataBase(db1)
		if err != nil {
			log.Fatal("open", err)
		}
		defer db1.Close()

		if err := db1.Ping(); err != nil {
			log.Fatal("ping", err)
		}
		server = handlerserver.New(nil, baseURL, host, datab)
	}

	// Запускаем сервер

	_, cancel := context.WithCancel(context.Background())

	go func() {
		log.Println("Запуск сервера на", host)
		if err := server.Start(); err != nil {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	fmt.Println("Остановка сервера...")
	cancel()

	// Даём время для завершения активных соединений
	time.Sleep(2 * time.Second)

	// Сохранение данных перед выходом
	if dboopen == "" {
		if err := db.WriteLocal(storagePath); err != nil {
			fmt.Printf("Ошибка сохранения данных: %v\n", err)
		} else {
			fmt.Println("Данные успешно сохранены")
		}
	}
}
