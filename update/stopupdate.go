package update

// TO DO: StopLat and StopLon to doubles or something similar

type STAPIResponse struct {
	Status   string     `json:"status"`
	Entities STEntities `json:"response"`
	Error    string     `json:"error"`
}

type STEntities []STEntity

type STEntity struct {
	StopID   string `json:"stop_id"`
	StopCode string `json:"stop_code"`
	StopName string `json:"stop_name"`
	StopLat  string `json:"stop_lat"`
	StopLon  string `json:"stop_lon"`
}
