package nats

import (
	"time"
)

// JiraIssueReassignRequest là cấu trúc request gửi đến Jira để reassign issue
type JiraIssueReassignRequest struct {
	JiraKey      string     `json:"jira_key"`    // Jira issue key
	NodeId       string     `json:"node_id"`     // Node ID trong hệ thống
	OldUserId    int32      `json:"old_user_id"` // ID người dùng cũ
	NewUserId    int32      `json:"new_user_id"` // ID người dùng mới
	LastSyncedAt *time.Time `json:"last_synced_at,omitempty"`
}

// JiraIssueReassignResponse là cấu trúc response nhận về từ Jira
type JiraIssueReassignResponse struct {
	Success      bool      `json:"success"`       // Trạng thái thành công hay không
	JiraKey      string    `json:"jira_key"`      // Jira issue key
	NodeId       string    `json:"node_id"`       // Node ID trong hệ thống
	OldUserId    int32     `json:"old_user_id"`   // ID người dùng cũ
	NewUserId    int32     `json:"new_user_id"`   // ID người dùng mới
	Timestamp    time.Time `json:"timestamp"`     // Thời điểm xử lý
	ErrorMessage *string   `json:"error_message"` // Thông báo lỗi nếu có
}
