package input

import (
	"fmt"
	"log"

	"github.com/autlamps/delay-backend-transformation/update"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (is *InService) SttIn(entities update.STTEntities) {

	tx, err := is.Db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(pq.CopyIn("stop_times", "stoptime_id", "trip_id", "arrival_time", "departure_time", "stop_id", "stop_sequence"))

	for i := 0; i < len(entities); i++ {
		gtfs_trip_id := entities[i].TripID
		gtfs_stop_id := entities[i].StopID
		stop_id := is.StopMap[gtfs_stop_id]
		trip_id := is.TripMap[gtfs_trip_id]

		arrival_time := entities[i].ArrivalTime
		departure_time := entities[i].DepatureTime
		stop_sequence := entities[i].StopSequence
		stoptime_id, err := uuid.NewRandom()
		if err != nil {
			log.Fatal(err)
		}
		gtfs_stop_time := fmt.Sprintf("%v%v", stop_id, trip_id)
		is.StopTimeMap[gtfs_stop_time] = stoptime_id

		gtfsID, itemuuid := is.toGTFSMap(is.StopTimeMap, gtfs_stop_time)
		is.GTFSStopTimeMap[gtfsID] = itemuuid

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
}
