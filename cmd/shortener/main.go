package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ippolid/shortLink/config"
	"github.com/Ippolid/shortLink/internal/app"
	"github.com/Ippolid/shortLink/internal/app/handlerserver"
	"github.com/Ippolid/shortLink/internal/logger"
)

func main() {
	// Получаем параметры конфигурации
	host, baseURL, storagePath := config.ParseFlags()

	// Инициализация логгера
	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}

	// Создаём базу данных
	db := app.NewDbase()

	// Загружаем данные из файла, если файл указан
	if err := db.ReadLocal(storagePath); err != nil {
		log.Fatalf("Ошибка загрузки данных: %v", err)
	}

	// Запускаем сервер
	server := handlerserver.New(&db, baseURL, host)
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
	if err := db.WriteLocal(storagePath); err != nil {
		fmt.Printf("Ошибка сохранения данных: %v\n", err)
	} else {
		fmt.Println("Данные успешно сохранены")
	}
}
