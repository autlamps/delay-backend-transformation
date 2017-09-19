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

// Needs Work
func Backup(db *sql.DB) {
	db.Exec("DROP SCHEMA backup CASCADE;")
	db.Exec("ALTER SCHEMA public RENAME TO backup;")
	db.Exec("CREATE SCHEMA public;")
}
