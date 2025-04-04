package queryparams

type RequestTaskQueryParam struct {
	Page     int
	PageSize int
}

type RequestTaskProjectQueryParam struct {
	Page         int
	PageSize     int
	WorkflowType string
	Status       string
	Type         string
	ProjectKey   string
}
