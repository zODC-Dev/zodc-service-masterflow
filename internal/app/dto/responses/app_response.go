package responses

type Response struct {
	Error    interface{} `json:"error"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
	Metadata interface{} `json:"metadata"`
}
