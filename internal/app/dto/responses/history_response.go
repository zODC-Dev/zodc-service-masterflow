package responses

import (
	"time"

	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
)

type HistoryNodeResponse struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Title string `json:"title"`
}

type HistoryResponse struct {
	ID        int32               `json:"id"`
	CreatedAt time.Time           `json:"createdAt"`
	Assignee  types.Assignee      `json:"assignee"`
	Type      string              `json:"type"`
	Node      HistoryNodeResponse `json:"node"`
	From      interface{}         `json:"from"`
	To        interface{}         `json:"to"`
}
