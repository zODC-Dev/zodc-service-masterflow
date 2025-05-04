package constants

const (
	FormTemplateIDJiraSystemForm      int32 = 1
	FormTemplateIDEditProfileForm     int32 = 2
	FormTemplateIDPerformanceEvaluate int32 = 7
)

type NodeFormPermission string

const (
	NodeFormPermissionInput  NodeFormPermission = "INPUT"
	NodeFormPermissionView   NodeFormPermission = "VIEW"
	NodeFormPermissionEdit   NodeFormPermission = "EDIT"
	NodeFormPermissionHidden NodeFormPermission = "HIDDEN"
)

type FormTemplateFieldType string

const (
	FormTemplateFieldTypeAttachment FormTemplateFieldType = "ATTACHMENT"
)
