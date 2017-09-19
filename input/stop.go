package input

import (
	"database/sql"
	"fmt"
	"github.com/autlamps/delay-backend-transformation/update"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
)

func StIn(entities update.STEntities, db *sql.DB, m map[string]uuid.UUID) {

	tx, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(pq.CopyIn("stops", "stop_id", "stop_code", "stop_name", "stop_lat", "stop_lon"))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(entities); i++ {
		stop_code := entities[i].StopCode
		stop_name := entities[i].StopName
		stop_lat := entities[i].StopLat
		stop_lon := entities[i].StopLon
		stop_id, err := uuid.NewRandom()
		if err != nil {
			log.Fatal(err)
		}
		m[entities[i].StopID] = stop_id
		_, err = stmt.Exec(stop_id, stop_code, stop_name, stop_lat, stop_lon)
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
	fmt.Println("Done Stops")
}