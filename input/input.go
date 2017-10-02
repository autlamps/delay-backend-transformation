package input

import (
	"database/sql"

	"strings"

	"github.com/google/uuid"
)

type InService struct {
	Db              *sql.DB
	Version         string
	GTFSAgencyMap   map[string]uuid.UUID
	GTFSServiceMap  map[string]uuid.UUID
	GTFSRouteMap    map[string]uuid.UUID
	GTFSTripMap     map[string]uuid.UUID
	GTFSStopMap     map[string]uuid.UUID
	GTFSStopTimeMap map[string]uuid.UUID
	AgencyMap       map[string]uuid.UUID
	ServiceMap      map[string]uuid.UUID
	RouteMap        map[string]uuid.UUID
	TripMap         map[string]uuid.UUID
	StopMap         map[string]uuid.UUID
	StopTimeMap     map[string]uuid.UUID
}

func (in *InService) Init() {
	in.AgencyMap = make(map[string]uuid.UUID)
	in.ServiceMap = make(map[string]uuid.UUID)
	in.RouteMap = make(map[string]uuid.UUID)
	in.TripMap = make(map[string]uuid.UUID)
	in.StopMap = make(map[string]uuid.UUID)
	in.StopTimeMap = make(map[string]uuid.UUID)

	in.GTFSAgencyMap = make(map[string]uuid.UUID)
	in.GTFSServiceMap = make(map[string]uuid.UUID)
	in.GTFSRouteMap = make(map[string]uuid.UUID)
	in.GTFSTripMap = make(map[string]uuid.UUID)
	in.GTFSStopMap = make(map[string]uuid.UUID)
	in.GTFSStopTimeMap = make(map[string]uuid.UUID)
}

func (in *InService) toGTFSMap(inMap map[string]uuid.UUID, gtfs string) (string, uuid.UUID) {
	itemuuid := inMap[gtfs]

	gtfsWversion := in.removeVersion(in.Version)
	return gtfsWversion, itemuuid
}

func (in *InService) removeVersion(str string) string {
	toReturn := strings.Replace(str, in.Version, "", -1)

	return toReturn
}
