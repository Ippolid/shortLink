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
CREATE TABLE IF NOT EXISTS public.shorty (
    id   TEXT PRIMARY KEY,
    link TEXT NOT NULL
);
`
	Insert = `
INSERT INTO public.shorty (id, link)
VALUES ($1, $2)
ON CONFLICT (id)
DO UPDATE SET link = EXCLUDED.link
`
	Get = `SELECT link FROM shorty WHERE id=$1`
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

func (data *DataBase) InsertLink(id, link string) error {
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
	_, err = tx.Exec(Insert, id, link)
	if err != nil {
		return fmt.Errorf("ошибка вставки данных: %v", err)
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
