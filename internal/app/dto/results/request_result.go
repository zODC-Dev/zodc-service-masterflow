package results

type RequestOverviewResult struct {
	MyRequests  int32 `json:"myrequest"`
	InProcess   int32 `json:"in_process"`
	Completed   int32 `json:"completed"`
	AllRequests int32 `json:"all_request"`
}
