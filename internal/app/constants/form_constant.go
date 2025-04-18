package constants

const (
	FormTemplateIDJiraSystemForm int32 = 1
)

type NodeFormPermission string

const (
	NodeFormPermissionInput NodeFormPermission = "INPUT"
	NodeFormPermissionView  NodeFormPermission = "VIEW"
	NodeFormPermissionEdit  NodeFormPermission = "EDIT"
)

type FormTemplateFieldType string

const (
	FormTemplateFieldTypeAttachment FormTemplateFieldType = "ATTACHMENT"
)
