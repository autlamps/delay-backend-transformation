package update

type CAAPIResponse struct {
	Status   string     `json:"status"`
	Entities CAEntities `json:"response"`
	Error    error      `json:"error"`
}

type CAEntities []CAEntity

type CAEntity struct {
	ServiceID string `json:"service_id"`
	Monday    int    `json:"monday"`
	Tuesday   int    `json:"tuesday"`
	Wednesday int    `json:"wednesday"`
	Thursday  int    `json:"thursday"`
	Friday    int    `json:"friday"`
	Saturday  int    `json:"saturday"`
	Sunday    int    `json:"Wednesday"`
}
