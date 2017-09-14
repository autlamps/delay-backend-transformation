package input

import (
	"database/sql"
	"fmt"
	"github.com/autlamps/delay-backend-transformation/update"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
)

func TrIn(entities update.TREntities, db *sql.DB, m map[string]uuid.UUID) {

	tx, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(pq.CopyIn("trips", "gtfs_trip_id", "trip_id", "route_id", "trip_headsign", "service_id"))

	for i := 0; i < len(entities); i++ {
		route_id := m[entities[i].RouteID]
		service_id := m[entities[i].ServiceID]
		gtfs_trip_id := entities[i].TripID
		trip_headsign := entities[i].TripHeadSign
		trip_id, err := uuid.NewRandom()

		if err != nil {
			log.Fatal(err)
		}

		m[gtfs_trip_id] = trip_id

		_, err = stmt.Exec(gtfs_trip_id, trip_id, route_id, trip_headsign, service_id)
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
	fmt.Println("Done Trips")
}
