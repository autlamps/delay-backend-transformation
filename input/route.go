package input

import (
	"log"

	"github.com/autlamps/delay-backend-transformation/update"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (is *InService) RoIn(entities update.ROEntities) {

	tx, err := is.Db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(pq.CopyIn("routes", "gtfs_route_id", "route_id", "agency_id", "route_short_name", "route_long_name", "route_type"))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(entities); i++ {
		gtfs_route_id := entities[i].RouteID
		gtfs_agency_id := entities[i].AgencyID
		route_short_name := entities[i].RouteSName
		route_long_name := entities[i].RouteLName
		route_type := entities[i].RouteType
		agency_id := is.AgencyMap[gtfs_agency_id]

		route_id, err := uuid.NewRandom()

		if err != nil {
			log.Fatal(err)
		}

		is.RouteMap[gtfs_route_id] = route_id
		gtfsID, iteamuuid := is.toGTFSMap(is.RouteMap, gtfs_route_id)
		is.GTFSRouteMap[gtfsID] = iteamuuid

		_, err = stmt.Exec(gtfs_route_id, route_id, agency_id, route_short_name, route_long_name, route_type)
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
