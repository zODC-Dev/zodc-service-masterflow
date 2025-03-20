package responses

type Response struct {
	Error    interface{} `json:"error"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
	Metadata interface{} `json:"metadata"`
}

type Paginate[T any] struct {
	Items      T
	Total      int
	Page       int
	PageSize   int
	TotalPages int
}
