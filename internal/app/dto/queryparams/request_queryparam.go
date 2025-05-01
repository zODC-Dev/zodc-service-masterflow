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

type RequestSubRequestQueryParam struct {
	Page     int
	PageSize int
}

type RequestTaskCount struct {
	WorkflowType string
	Type         string
	ProjectKey   string
}

type RequestMidSprintReportQueryParam struct {
	StartTime  string
	EndTime    string
	SprintId   string
	ProjectKey string
}
