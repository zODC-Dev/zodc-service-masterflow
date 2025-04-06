package responses

import (
	"time"

	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
)

type RequestCountResponse struct {
	MyRequests  int32 `json:"myrequest"`
	InProcess   int32 `json:"in_process"`
	Completed   int32 `json:"completed"`
	AllRequests int32 `json:"all_request"`
}

type CurrentTaskResponse struct {
	Title        string           `json:"title"`
	UpdatedAt    time.Time        `json:"updatedAt"`
	Participants []types.Assignee `json:"participants"`
}

type RequestResponse struct {
	Id           int32                 `json:"id"`
	Key          int32                 `json:"key"`
	Title        string                `json:"title"`
	Status       string                `json:"status"`
	ParentKey    int32                 `json:"parentKey"`
	Progress     int32                 `json:"progress"`
	CurrentTasks []CurrentTaskResponse `json:"currentTasks,omitempty"`
	SprintId     int32                 `json:"sprintId"`
	StartedAt    *time.Time            `json:"startedAt"`
	CompletedAt  *time.Time            `json:"completedAt"`
	CanceledAt   *time.Time            `json:"canceledAt"`
	TerminatedAt *time.Time            `json:"terminatedAt"`
}

type RequestDetailResponse struct {
	RequestResponse
	ParentRequest *RequestResponse  `json:"parentRequest"`
	ChildRequests []RequestResponse `json:"childRequests"`
	RequestedBy   types.Assignee    `json:"requestedBy"`
	Participants  []types.Assignee  `json:"participants"`
	Workflow      WorkflowResponse  `json:"workflow"`
}

type RequestTaskResponse struct {
	Id               string         `json:"id"`
	Key              string         `json:"key"`
	Title            string         `json:"title"`
	Type             string         `json:"type"`
	RequestID        int32          `json:"requestId"`
	RequestTitle     string         `json:"requestTitle"`
	RequestProgress  float32        `json:"requestProgress"`
	Assignee         types.Assignee `json:"assignee"`
	Status           string         `json:"status"`
	PlannedStartTime *time.Time     `json:"plannedStartTime"`
	PlannedEndTime   *time.Time     `json:"plannedEndTime"`
	ActualStartTime  *time.Time     `json:"actualStartTime"`
	ActualEndTime    *time.Time     `json:"actualEndTime"`
	EstimatePoint    *int32         `json:"estimatePoint"`
	IsCurrent        bool           `json:"isCurrent"`
}

type TaskCountResponse struct {
	OverdueCount   int32 `json:"overdueCount"`
	TotalCount     int32 `json:"totalCount"`
	CompletedCount int32 `json:"completedCount"`
	TodoCount      int32 `json:"todoCount"`
}

type RequestOverviewResponse struct {
	WorkflowDetailResponse
	Progress float32 `json:"progress"`
}
