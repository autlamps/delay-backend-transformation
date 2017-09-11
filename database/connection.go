package database

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateCon() *sql.DB {
	db, err := sql.Open("postgres", "postgresql://postgres:postgres@127.0.0.1/gtfs?sslmode=disable")

	if err != nil {
		fmt.Println(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
