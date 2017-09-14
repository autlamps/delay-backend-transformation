package input

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/autlamps/delay-backend-transformation/update"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func SttIn(entities update.STTEntities, db *sql.DB, m map[string]uuid.UUID) {

	tx, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(pq.CopyIn("stop_times", "stoptime_id", "trip_id", "arrival_time", "departure_time", "stop_id", "stop_sequence"))

	for i := 0; i < len(entities); i++ {
		stop_id := m[entities[i].StopID]
		trip_id := m[entities[i].TripID]
		arrival_time := entities[i].ArrivalTime
		departure_time := entities[i].DepatureTime
		stop_sequence := entities[i].StopSequence
		stoptime_id, err := uuid.NewRandom()
		if err != nil {
			log.Fatal(err)
		}
		gtfs_stop_time := fmt.Sprint("%v%v", stop_id, trip_id)
		m[gtfs_stop_time] = stoptime_id

		_, err = stmt.Exec(stoptime_id, trip_id, arrival_time, departure_time, stop_id, stop_sequence)
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

	fmt.Println("Done Stoptime")
}
