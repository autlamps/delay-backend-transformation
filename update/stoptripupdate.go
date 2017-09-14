package update

import "time"

type STTAPIResponse struct {
	Status   string      `json:"status"`
	Entities STTEntities `json:"response"`
	Error    error       `json:"error"`
}

type STTEntities []STTEntity

type STTEntity struct {
	TripID       string
	ArrivalTime  time.Time
	DepatureTime time.Time
	StopID       string
	StopSequence int
}
