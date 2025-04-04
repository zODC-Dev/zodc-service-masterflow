package queryparams

type WorkflowQueryParam struct {
	CategoryID     string
	Type           string
	Search         string
	ProjectKey     string
	HasSubWorkflow string
	IsArchived     string
}

type RequestQueryParam struct {
	Search       string
	Page         int
	PageSize     int
	ProjectKey   string
	Status       string
	SprintID     string
	WorkflowType string
}
