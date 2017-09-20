package input

import (
	"fmt"
	"log"

	"github.com/autlamps/delay-backend-transformation/update"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (is *InService) AgIn(entities update.AGEntities) {

	tx, err := is.Db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(pq.CopyIn("agency", "gtfs_agency_id", "agency_name", "agency_id"))

	for i := 0; i < len(entities); i++ {
		gtfs_agency_id := entities[i].AgencyID
		agency_name := entities[i].AgencyName
		agency_id, err := uuid.NewRandom()
		if err != nil {
			fmt.Println("No new UUID")
			log.Fatal(err.Error())
		}

		is.AgencyMap[gtfs_agency_id] = agency_id

		_, err = stmt.Exec(gtfs_agency_id, agency_name, agency_id)
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
