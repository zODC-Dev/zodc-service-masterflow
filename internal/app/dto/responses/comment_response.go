package responses

import (
	"time"

	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
)

type CommentResponse struct {
	Content   string         `json:"content"`
	Assignee  types.Assignee `json:"assignee"`
	CreatedAt time.Time      `json:"createdAt"`
}
