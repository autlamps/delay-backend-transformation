package update

type VRAPIResponse struct {
	Status   string  `json:"status"`
	ATVerDet []Version `json:"response"`
	Error    error   `json:"error"`
}

type Versions []Version

type Version struct {
	ATVersion string `json:"version"`
	StartDate string `json:"startdate"`
	EndDate   string `json:"enddate"`
}

type VerResp struct {
	Version Version
	Error   error
}
