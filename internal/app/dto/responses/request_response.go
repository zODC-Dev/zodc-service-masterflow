package responses

import (
	"time"

	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
)

type RequestOverviewResponse struct {
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
	Title        int32                 `json:"title"`
	ParentKey    int32                 `json:"parentKey"`
	Progress     int32                 `json:"progress"`
	CurrentTasks []CurrentTaskResponse `json:"currenTasks"`
	SprintId     int32                 `json:"sprintId"`
	StartedAt    *time.Time            `json:"startedAt"`
	CompletedAt  *time.Time            `json:"completedAt"`
}
