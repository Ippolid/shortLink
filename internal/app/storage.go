package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Ippolid/shortLink/internal/models"

	//"io"
	//"io/ioutil"
	"os"
)

func (s *Dbase) ReadLocal(path string) error {
	if path == "" {
		return fmt.Errorf("path is empty")
	}
	var db = make([]models.LocalLink, 0)

	// Открываем файл
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	// Декодируем JSON-данные
	if err := decoder.Decode(&db); err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	for _, v := range db {
		s.Data[v.ID] = v.URL
	}

	return nil
}

func (s *Dbase) WriteLocal(path string) error {
	if path == "" {
		return fmt.Errorf("path is empty")
	}

	// Преобразуем данные в срез структур LocalLink
	var db []models.LocalLink
	for id, url := range s.Data {
		db = append(db, models.LocalLink{ID: id, URL: url})
	}

	fmt.Println(db)

	// Открываем файл для записи
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)

	// Кодируем данные в JSON и записываем в файл
	if err := encoder.Encode(&db); err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	// Обязательно сбрасываем буфер в файл
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %v", err)
	}

	return nil
}
