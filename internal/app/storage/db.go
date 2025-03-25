package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"sync"
)

// DataBase — обёртка для работы с БД
type DataBase struct {
	db *sql.DB
}

// Const - создание констант
const (
	//CrearTable — запрос на создание таблицы
	CrearTable = `
	    CREATE TABLE IF NOT EXISTS shorty (
	        id      TEXT PRIMARY KEY,
	        link    TEXT NOT NULL UNIQUE,
	        user_id TEXT default '',
	        deleted BOOLEAN default false
	    );`

	//Insert — запрос на вставку данных
	Insert = `
	    INSERT INTO shorty (id, link, user_id)
	    VALUES ($1, $2, $3)
	    ON CONFLICT DO NOTHING;
	    `
	//Get — запрос на получение данных
	Get = `SELECT link,deleted FROM shorty WHERE id=$1`
	//GetLinksByUserID — запрос на получение данных по пользователю
	GetLinksByUserID = `SELECT link,deleted FROM shorty WHERE user_id=$1`
	//DelLink — запрос на удаление данных
	DelLink = `UPDATE shorty SET deleted = TRUE
	WHERE id = $1 AND user_id = $2`
)

// NewDBWrapper — конструктор для обёртки
func NewDataBase(db *sql.DB) (*DataBase, error) {
	_, err := db.Exec(CrearTable)
	if err != nil {
		return nil, err
	}
	return &DataBase{db: db}, nil
}

// Connect — подключение к БД
func Connect(op string) (*sql.DB, error) {
	ps := op

	db, err := sql.Open("pgx", ps)
	fmt.Println(db, err)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка пинга БД: %v", err)
	}

	return db, nil
}

// Ping — проверка доступности БД
func (data *DataBase) Ping() (bool, error) {
	err := data.db.Ping()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	return true, nil
}

// InsertLink — вставка данных
func (data *DataBase) InsertLink(id, link, user string) error {
	tx, err := data.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %v", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Ошибка отката транзакции: %v", rollbackErr)
			}
		}
	}()

	result, err := tx.Exec(Insert, id, link, user)
	if err != nil {
		return fmt.Errorf("link exists %v", err)
	}

	// Проверяем, была ли выполнена вставка
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества затронутых строк: %v", err)
	}

	if rowsAffected == 0 {
		// Если вставка не произошла из-за ON CONFLICT DO NOTHING,
		// возвращаем специальную ошибку, которую потом можно обработать
		return fmt.Errorf("link exists %v", err)
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %v", err)
	}

	log.Println("Данные успешно вставлены")
	return nil
}

// GetLink — получение данных
func (data *DataBase) GetLink(id string) (string, bool, error) {
	var link string
	var deleted bool
	err := data.db.QueryRow(Get, id).Scan(&link, &deleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, fmt.Errorf("запись не найдена")
		}
		return "", false, fmt.Errorf("ошибка получения данных: %v", err)
	}
	return link, deleted, nil
}

// GetLinksByUserID — получение данных по пользователю
func (data *DataBase) GetLinksByUserID(userID string) ([]string, error) {
	rows, err := data.db.Query(GetLinksByUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	var deleted bool
	defer rows.Close()

	var links []string
	for rows.Next() {
		var link string
		if err := rows.Scan(&link, &deleted); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}
		if !deleted {
			links = append(links, link)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка обработки строк: %v", err)
	}

	return links, nil
}

// Dellink — удаление данных
func (data *DataBase) Dellink(ids []string, userID string) error {
	buf := make(chan string, 100)
	var wg sync.WaitGroup

	// Запускаем горутины для обработки удаления
	for i := 0; i < len(ids); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range buf {
				_, err := data.db.Exec(DelLink, id, userID)
				if err != nil {
					log.Printf("Ошибка при удалении ссылки %s: %v", id, err)
				} // Небольшая задержка для снижения нагрузки на БД
			}
		}()
	}

	// Отправляем идентификаторы в буфер
	for _, id := range ids {
		buf <- id
	}
	close(buf)

	// Ждём завершения всех горутин
	wg.Wait()
	return nil
}
