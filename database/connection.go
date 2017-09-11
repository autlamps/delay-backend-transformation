package database

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateCon(DB_URL string) *sql.DB {
	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		fmt.Println(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
