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
}

type JiraFormDetailResponse struct {
	Template FormTemplateFindAll         `json:"template"`
	Fields   []FormTemplateFieldsFindAll `json:"fields"`
	Data     []NodeFormDataResponse      `json:"data"`
}

type TaskDetail struct {
	RequestTaskResponse
	UpdatedAt        time.Time      `json:"updatedAt"`
	RequestRequestBy types.Assignee `json:"requestRequestBy"`
	IsApproval       bool           `json:"isApproval"`
}
