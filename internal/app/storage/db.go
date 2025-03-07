package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

type DataBase struct {
	db *sql.DB
}

const (
	CrearTable = `
	    CREATE TABLE IF NOT EXISTS shorty (
	        id      TEXT PRIMARY KEY,
	        link    TEXT NOT NULL UNIQUE,
	        user_id TEXT default '123123'
	    );
	    `
	Insert = `
	    INSERT INTO shorty (id, link, user_id)
	    VALUES ($1, $2, $3)
	    ON CONFLICT DO NOTHING;
	    `
	Get              = `SELECT link FROM shorty WHERE id=$1`
	GetLinksByUserID = `SELECT link FROM shorty WHERE user_id=$1`
)

// NewDBWrapper — конструктор для обёртки
func NewDataBase(db *sql.DB) (*DataBase, error) {
	_, err := db.Exec(CrearTable)
	if err != nil {
		return nil, err
	}
	return &DataBase{db: db}, nil
}

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

func (data *DataBase) Ping() (bool, error) {
	err := data.db.Ping()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	return true, nil
}

func (data *DataBase) InsertLink(id, link, user string) error {
	//tx, err := data.db.Begin()
	//if err != nil {
	//	return fmt.Errorf("ошибка начала транзакции: %v", err)
	//}
	//defer func() {
	//	if err != nil {
	//		if rollbackErr := tx.Rollback(); rollbackErr != nil {
	//			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
	//		}
	//	}
	//}()
	//res, err1 := tx.Exec(Insert, id, link)
	//
	//fmt.Println(res, err1)
	//
	//if err != nil {
	//	return fmt.Errorf("ошибка вставки данных: %v", err)
	//}
	//
	//// Фиксируем транзакцию
	//if err = tx.Commit(); err != nil {
	//	return fmt.Errorf("ошибка подтверждения транзакции: %v", err)
	//}
	//
	//log.Println("Данные успешно вставлены")
	//return nil
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

func (data *DataBase) GetLink(id string) (string, error) {
	var link string
	err := data.db.QueryRow(Get, id).Scan(&link)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("запись не найдена")
		}
		return "", fmt.Errorf("ошибка получения данных: %v", err)
	}
	return link, nil
}

func (data *DataBase) GetLinksByUserID(userID string) ([]string, error) {
	rows, err := data.db.Query(GetLinksByUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var links []string
	for rows.Next() {
		var link string
		if err := rows.Scan(&link); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка обработки строк: %v", err)
	}

	return links, nil
}
