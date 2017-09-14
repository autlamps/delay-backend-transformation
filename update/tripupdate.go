package update

type TRAPIResponse struct {
	Status   string     `json:"status"`
	Entities TREntities `json:"response"`
	Error    error      `json:"error"`
}

type TREntities []TREntity

type TREntity struct {
	RouteID      string `json:"route_id"`
	ServiceID    string `json:"service_id"`
	TripID       string `json:"trip_id"`
	TripHeadSign string `json:"trip_headsign"`
	TripSName    string `json:"trip_short_name"`
}
