// package app
//
// import (
//
//	"bufio"
//	"encoding/json"
//	"fmt"
//	"github.com/Ippolid/shortLink/internal/models"
//
//	//"io"
//	//"io/ioutil"
//	"os"
//
// )
//
//	func (s *Dbase) ReadLocal(path string) error {
//		if path == "" {
//			return fmt.Errorf("path is empty")
//		}
//		var db = make([]models.LocalLink, 0)
//
//		// Открываем файл
//		file, err := os.Open(path)
//		if err != nil {
//			return fmt.Errorf("error opening file: %v", err)
//		}
//		defer file.Close()
//
//		reader := bufio.NewReader(file)
//		decoder := json.NewDecoder(reader)
//
//		// Декодируем JSON-данные
//		if err := decoder.Decode(&db); err != nil {
//			return fmt.Errorf("error decoding JSON: %v", err)
//		}
//
//		for _, v := range db {
//			s.Data[v.ID] = v.URL
//		}
//
//		return nil
//	}
//
//	func (s *Dbase) WriteLocal(path string) error {
//		if path == "" {
//			return fmt.Errorf("path is empty")
//		}
//
//		// Преобразуем данные в срез структур LocalLink
//		var db []models.LocalLink
//		for id, url := range s.Data {
//			db = append(db, models.LocalLink{ID: id, URL: url})
//		}
//
//		fmt.Println(db)
//
//		// Открываем файл для записи
//		file, err := os.Create(path)
//		if err != nil {
//			return fmt.Errorf("error creating file: %v", err)
//		}
//		defer file.Close()
//
//		writer := bufio.NewWriter(file)
//		encoder := json.NewEncoder(writer)
//
//		// Кодируем данные в JSON и записываем в файл
//		if err := encoder.Encode(&db); err != nil {
//			return fmt.Errorf("error encoding JSON: %v", err)
//		}
//
//		// Обязательно сбрасываем буфер в файл
//		if err := writer.Flush(); err != nil {
//			return fmt.Errorf("error flushing writer: %v", err)
//		}
//
//		return nil
//	}
package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Ippolid/shortLink/internal/models"
)

// ReadLocal загружает данные из файла JSON
func (s *Dbase) ReadLocal(path string) error {
	if path == "" {
		fmt.Println("Файл хранилища отключен, загрузка пропущена")
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Если файла нет, создаём пустой
			fmt.Println("Файл не найден, создаётся новый:", path)
			return s.WriteLocal(path)
		}
		return fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(bufio.NewReader(file))
	var db []models.LocalLink
	if err := decoder.Decode(&db); err != nil {
		return fmt.Errorf("ошибка декодирования JSON: %v", err)
	}

	// Загружаем данные в память
	for _, entry := range db {
		s.Data[entry.ID] = entry.URL
	}

	fmt.Println("Данные успешно загружены из файла")
	return nil
}

// WriteLocal сохраняет данные в файл JSON
func (s *Dbase) WriteLocal(path string) error {
	if path == "" {
		fmt.Println("Файл хранилища отключен, сохранение пропущено")
		return nil
	}

	// Преобразуем данные в срез структур LocalLink
	var db []models.LocalLink
	for id, url := range s.Data {
		db = append(db, models.LocalLink{ID: id, URL: url})
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)

	if err := encoder.Encode(db); err != nil {
		return fmt.Errorf("ошибка кодирования JSON: %v", err)
	}

	// Обязательно сбрасываем буфер в файл
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("ошибка записи в файл: %v", err)
	}

	fmt.Println("Данные успешно сохранены в файл")
	return nil
}
