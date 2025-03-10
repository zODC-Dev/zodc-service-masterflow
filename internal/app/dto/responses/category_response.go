package responses

type CategoryFindAll struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Key      string `json:"key"`
	IsActive bool   `json:"isActive"`
}
