package static

import (
	"database/sql"
	"github.com/google/uuid"
)

type Agency struct {
	ID         uuid.UUID
	GTFSID     string
	AgencyName string
}

type AgencyStore interface {
	GetAgencybyID(id uuid.UUID) (Agency, error)
}

type AgencyService struct {
	db *sql.DB
}

func AgencyServiceInit(db *sql.DB) *AgencyService {
	return &AgencyService{db: db}
}

func (as *AgencyService) GetAgencyByID(id uuid.UUID) (Agency, error) {
	a := Agency{}

	row := as.db.QueryRow("SELECT agency_id, gtfs_agency_id, agency_name FROM agency WHERE agency_id = $1", id)
	err := row.Scan(&a.ID, &a.GTFSID, &a.AgencyName)

	if err != nil {
		return a, err
	}

	return a, nil
}
