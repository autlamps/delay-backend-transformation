package input

import (
	"database/sql"

	"github.com/google/uuid"
)

type InService struct {
	Db          *sql.DB
	AgencyMap   map[string]uuid.UUID
	ServiceMap  map[string]uuid.UUID
	RouteMap    map[string]uuid.UUID
	TripMap     map[string]uuid.UUID
	StopMap     map[string]uuid.UUID
	StopTimeMap map[string]uuid.UUID
}

func (in *InService) Init() {
	in.AgencyMap = make(map[string]uuid.UUID)
	in.ServiceMap = make(map[string]uuid.UUID)
	in.RouteMap = make(map[string]uuid.UUID)
	in.TripMap = make(map[string]uuid.UUID)
	in.StopMap = make(map[string]uuid.UUID)
	in.StopTimeMap = make(map[string]uuid.UUID)
}
