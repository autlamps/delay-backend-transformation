package input

import (
	"github.com/autlamps/delay-backend-transformation/update"
	"database/sql"
	"github.com/google/uuid"
	"log"
	"fmt"
)

func RoIn(entities update.ROEntities, db *sql.DB, m map[string]uuid.UUID) {
	for i := 0; i < len(entities); i++ {
		gtfs_route_id := entities[i].RouteID
		gtfs_agency_id := entities[i].AgencyID
		route_short_name := entities[i].RouteSName
		route_long_name := entities[i].RouteLName
		route_type := entities[i].RouteType
		agency_id := m[gtfs_agency_id]
		route_id, err := uuid.NewRandom()

		if err != nil {
			log.Fatal(err)
		}

		m[gtfs_route_id] = route_id

		db.Exec("INSERT INTO routes (gtfs_route_id, route_id, agency_id, route_short_name, route_long_name, route_type) VALUES ($1, $2, $3, $4, $5, $6);", gtfs_route_id,route_id,agency_id, route_short_name, route_long_name, route_type)

	}
	fmt.Println("Done Routes")
}
