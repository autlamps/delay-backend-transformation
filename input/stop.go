package input

import (
	"log"

	"github.com/autlamps/delay-backend-transformation/update"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (is *InService) StIn(entities update.STEntities) {

	tx, err := is.Db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(pq.CopyIn("stops", "stop_id", "stop_code", "stop_name", "stop_lat", "stop_lon"))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(entities); i++ {
		gtfs_stop_id := entities[i].StopID
		stop_code := entities[i].StopCode
		stop_name := entities[i].StopName
		stop_lat := entities[i].StopLat
		stop_lon := entities[i].StopLon
		stop_id, err := uuid.NewRandom()
		if err != nil {
			log.Fatal(err)
		}

		is.StopMap[gtfs_stop_id] = stop_id

		gtfsID, itemuuid := is.toGTFSMap(is.StopMap, gtfs_stop_id)
		is.GTFSStopMap[gtfsID] = itemuuid

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
}
