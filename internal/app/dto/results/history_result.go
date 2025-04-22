package results

import "github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"

type HistoryResult struct {
	model.Histories
	Node model.Nodes
}
