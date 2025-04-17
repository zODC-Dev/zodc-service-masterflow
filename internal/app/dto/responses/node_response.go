package responses

import (
	"time"

	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
)

type NodeFormDetailResponse struct {
	Template    FormTemplateFindAll           `json:"template"`
	Fields      [][]FormTemplateFieldsFindAll `json:"fields"`
	Data        []NodeFormDataResponse        `json:"data"`
	DataId      string                        `json:"dataId"`
	IsSubmitted bool                          `json:"isSubmitted"`
	IsApproved  bool                          `json:"isApproved"`
	IsRejected  bool                          `json:"isRejected"`
}

type JiraFormDetailResponse struct {
	Template FormTemplateFindAll           `json:"template"`
	Fields   [][]FormTemplateFieldsFindAll `json:"fields"`
	Data     []NodeFormDataResponse        `json:"data"`
}

type TaskRelated struct {
	Key      string         `json:"key"`
	Title    string         `json:"title"`
	Type     string         `json:"type"`
	Status   string         `json:"status"`
	Assignee types.Assignee `json:"assignee"`
}

type TaskDetail struct {
	RequestTaskResponse
	UpdatedAt        time.Time      `json:"updatedAt"`
	RequestRequestBy types.Assignee `json:"requestRequestBy"`
	IsApproval       bool           `json:"isApproval"`

	SprintId      *int          `json:"sprintId"`
	EstimatePoint *int          `json:"estimatePoint"`
	Parent        *TaskRelated  `json:"parent"`
	Related       []TaskRelated `json:"related"`

	JiraLinkURL *string `json:"jiraLinkUrl"`
	ProjectKey  *string `json:"projectKey"`
}
