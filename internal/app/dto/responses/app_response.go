package responses

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Paginate[T any] struct {
	Items      T   `json:"items"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalPages int `json:"totalPages"`
}
