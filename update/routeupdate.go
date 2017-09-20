package update

//ROAPIResponse is the full response from the API
type ROAPIResponse struct {
	Status   string     `json:"status"`
	Entities ROEntities `json:"response"`
	Error    error      `json:"error"`
}

// ROEntities is an array of all ROEntity
type ROEntities []ROEntity

// ROEntity contains information about the individual Route
type ROEntity struct {
	RouteID    string `json:"route_id"`
	AgencyID   string `json:"agency_id"`
	RouteSName string `json:"route_short_name"`
	RouteLName string `json:"route_long_name"`
	RouteType  int    `json:"route_type"`
}

type ROReturn struct {
	Entities ROEntities
	Error    error
}
