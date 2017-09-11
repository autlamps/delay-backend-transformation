package input

import (
	"github.com/autlamps/delay-backend-transformation/update"
	"database/sql"
	"github.com/google/uuid"
	"fmt"
)

func StIn(entities update.STEntities, db *sql.DB, m map[string]uuid.UUID) {
	for i := 0; i < len(entities); i++ {
		stop_code := entities[i].StopCode
		stop_name := entities[i].StopName
		stop_lat := entities[i].StopLat
		stop_lon := entities[i].StopLon
		stop_id := entities[i].StopID

		db.Exec("INSERT INTO stops (stop_id, stop_code, stop_name, stop_lat, stop_lon) VALUES ($1, $2, $3, $4, $5);",stop_id, stop_code, stop_name, stop_lat, stop_lon)
	}
	fmt.Println("Done Stops")
}
