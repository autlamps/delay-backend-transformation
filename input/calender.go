package input

import (
	"database/sql"
	"fmt"
	"github.com/autlamps/delay-backend-transformation/update"
	"github.com/google/uuid"
	"log"
)

func CaIn(entities update.CAEntities, db *sql.DB, m map[string]uuid.UUID) {
	for i := 0; i < len(entities); i++ {
		gtfs_service_id := entities[i].ServiceID
		monday := boolcheck(entities[i].Monday)
		tuesday := boolcheck(entities[i].Tuesday)
		wednesday := boolcheck(entities[i].Wednesday)
		thursday := boolcheck(entities[i].Thursday)
		friday := boolcheck(entities[i].Friday)
		saturday := boolcheck(entities[i].Saturday)
		sunday := boolcheck(entities[i].Sunday)
		service_id, err := uuid.NewRandom()
		if err != nil {
			log.Fatal(err)
		}
		m[gtfs_service_id] = service_id

		db.Exec("INSERT INTO calendar (gtfs_service_id, service_id, monday, tuesday, wednesday, thursday, friday, saturday, sunday) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);", gtfs_service_id, service_id, monday, tuesday, wednesday, thursday, friday, saturday, sunday)
	}
	fmt.Println("Done Calender")
}

func boolcheck(check int) bool {
	if check == 0 {
		return false
	} else {
		return true
	}
}
