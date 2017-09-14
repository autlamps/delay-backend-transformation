package update

// AUAPIResponse is the full respsonse from the API
type AGAPIResponse struct {
	Status   string     `json:"status"`
	Entities AGEntities `json:"response"`
	Error    error      `json:"error"`
}

// AUEntites is simply a slice of AUEntity
type AGEntities []AGEntity

// AIEntity contains information about an individual Agency
type AGEntity struct {
	AgencyName string `json:"agency_name"`
	AgencyID   string `json:"agency_id"`
}
