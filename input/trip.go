package input

import (
	"github.com/autlamps/delay-backend-transformation/update"
	"database/sql"
	"github.com/google/uuid"
	"log"
	"fmt"
)

func TrIn( entities update.TREntities, db *sql.DB, m map[string]uuid.UUID) {
	for i := 0; i < len(entities); i++ {
		route_id := m[entities[i].RouteID]
		service_id := m[entities[i].ServiceID]
		gtfs_trip_id := entities[i].TripID
		trip_headsign := entities[i].TripHeadSign
		trip_short_name := entities[i].TripSName
		trip_id, err := uuid.NewRandom()

		if err != nil {
			log.Fatal(err)
		}

		m[gtfs_trip_id] = trip_id

		_, err = db.Exec("INSERT INTO trips (gtfs_trip_id, trip_id, route_id, trip_headsign, trip_short_name, service_id) VALUES ($1, $2, $3, $4, $5, $6);", gtfs_trip_id, trip_id, route_id, trip_headsign, trip_short_name, service_id)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Done Trips")
}
