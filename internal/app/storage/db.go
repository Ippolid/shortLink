package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

func Connect(op string) (*sql.DB, error) {
	ps := op

	db, err := sql.Open("pgx", ps)
	fmt.Println(db, err)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Ping(db *sql.DB) (bool, error) {
	err := db.Ping()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	return true, nil
}
