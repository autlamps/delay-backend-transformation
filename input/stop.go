package input

import (
	"database/sql"
	"fmt"
	"github.com/autlamps/delay-backend-transformation/update"
	"github.com/google/uuid"
	"log"
)

func StIn(entities update.STEntities, db *sql.DB, m map[string]uuid.UUID) {
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
		_, err = db.Exec("INSERT INTO stops (stop_id, stop_code, stop_name, stop_lat, stop_lon) VALUES ($1, $2, $3, $4, $5);", stop_id, stop_code, stop_name, stop_lat, stop_lon)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Done Stops")
}
