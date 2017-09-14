package input

import (
	"database/sql"
	"fmt"
	"github.com/autlamps/delay-backend-transformation/update"
	"github.com/google/uuid"
	"log"
	"github.com/lib/pq"
)

func CaIn(entities update.CAEntities, db *sql.DB, m map[string]uuid.UUID) {

	tx, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(pq.CopyIn("calendar", "gtfs_service_id", "service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"))

	for i := 0; i < len(entities); i++ {
		gtfs_service_id := entities[i].ServiceID
		monday := boolcheck(entities[i].Monday)
		tuesday := boolcheck(entities[i].Tuesday)
		wednesday := boolcheck(entities[i].Wednesday)
		thursday := boolcheck(entities[i].Thursday)
		friday := boolcheck(entities[i].Friday)
		saturday := boolcheck(entities[i].Saturday)
		sunday := boolcheck(entities[i].Sunday)
		service_id, err := uuid.NewRandom()
		if err != nil {
			log.Fatal(err)
		}
		m[gtfs_service_id] = service_id

		_, err = stmt.Exec(gtfs_service_id, service_id, monday, tuesday, wednesday, thursday, friday, saturday, sunday)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done Calender")
}

func boolcheck(check int) bool {
	if check == 0 {
		return false
	} else {
		return true
	}
}
