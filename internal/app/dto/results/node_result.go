package results

import (
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
)

type NodeResult struct {
	model.Nodes
	Request model.Requests
}
