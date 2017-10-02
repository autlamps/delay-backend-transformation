package input

import (
	"log"

	"github.com/autlamps/delay-backend-transformation/update"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (is *InService) TrIn(entities update.TREntities) {

	tx, err := is.Db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(pq.CopyIn("trips", "gtfs_trip_id", "trip_id", "route_id", "trip_headsign", "service_id"))

	for i := 0; i < len(entities); i++ {
		route_id := is.RouteMap[entities[i].RouteID]
		service_id := is.ServiceMap[entities[i].ServiceID]
		gtfs_trip_id := entities[i].TripID
		trip_headsign := entities[i].TripHeadSign
		trip_id, err := uuid.NewRandom()

		if err != nil {
			log.Fatal(err)
		}

		is.TripMap[gtfs_trip_id] = trip_id

		gtfsID, itemuuid := is.toGTFSMap(is.TripMap, gtfs_trip_id)
		is.GTFSTripMap[gtfsID] = itemuuid

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
}
