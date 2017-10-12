package input

import (
	"log"

	"github.com/autlamps/delay-backend-transformation/update"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (is *InService) CaIn(entities update.CAEntities) {

	tx, err := is.Db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	udte, err := tx.Prepare("UPDATE calendar SET gtfs_service_id = $1, monday = $2, tuesday = $3, wednesday = $4, thursday = $5, friday = $6, saturday = $7, sunday = $8 WHERE service_id = $9")
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(pq.CopyIn("calendar", "gtfs_service_id", "service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(entities); i++ {
		gtfs_service_id := entities[i].ServiceID
		monday := boolcheck(entities[i].Monday)
		tuesday := boolcheck(entities[i].Tuesday)
		wednesday := boolcheck(entities[i].Wednesday)
		thursday := boolcheck(entities[i].Thursday)
		friday := boolcheck(entities[i].Friday)
		saturday := boolcheck(entities[i].Saturday)
		sunday := boolcheck(entities[i].Sunday)

		itemuuid := is.GTFSServiceMap[is.removeVersion(entities[i].ServiceID)]

		if itemuuid != uuid.Nil {
			_, err = udte.Exec(gtfs_service_id, monday, tuesday, wednesday, thursday, friday, saturday, sunday, itemuuid)
			if err != nil {
				log.Fatal(err)
			}
		} else {

			service_id, err := uuid.NewRandom()
			if err != nil {
				log.Fatal(err)
			}
			is.ServiceMap[gtfs_service_id] = service_id
			gtfsID, itemuuid := is.toGTFSMap(is.ServiceMap, gtfs_service_id)
			is.GTFSServiceMap[gtfsID] = itemuuid

			_, err = stmt.Exec(gtfs_service_id, service_id, monday, tuesday, wednesday, thursday, friday, saturday, sunday)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = udte.Close()
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

func boolcheck(check int) bool {
	if check == 0 {
		return false
	} else {
		return true
	}
}
