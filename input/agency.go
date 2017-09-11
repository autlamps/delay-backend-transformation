package input

import (
	"database/sql"
	"fmt"
	"github.com/autlamps/delay-backend-transformation/update"
	"github.com/google/uuid"
	"log"
)

func AgIn(entities update.AGEntities, db *sql.DB, m map[string]uuid.UUID) {

	for i := 0; i < len(entities); i++ {
		gtfs_agency_id := entities[i].AgencyID
		agency_name := entities[i].AgencyName
		agency_id, err := uuid.NewRandom()
		if err != nil {
			fmt.Println("No new UUID")
			log.Fatal(err.Error())
		}
		m[gtfs_agency_id]	= agency_id

		db.Exec("INSERT INTO agency (gtfs_agency_id, agency_name, agency_id) VALUES ($1, $2, $3);", gtfs_agency_id, agency_name, agency_id)
	}
	fmt.Println("Done Agencies")
}


