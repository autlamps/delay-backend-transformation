package update

// AUAPIResponse is the full repsonse form the API
type AUAPIResponse struct {
	Status   string     `json:"status"`
	Entities AUEntities `json:"response"`
	Error    string     `json:"error"`
}

// AUEntites is simply a slice of AUEntity
type AUEntities []AUEntity

// AIEntity contains information about an individual Agency
type AUEntity struct {
	AgencyName string `json:"agency_name"`
	AgencyID   string `json:"agency_id"`
}
