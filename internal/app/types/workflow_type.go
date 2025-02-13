package types

import "github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/model"

type WorkflowType struct {
	model.Workflows
	Category    model.Categories
	Nodes       []model.Nodes
	Groups      []model.NodeGroups
	Connections []model.NodeConnections
}
