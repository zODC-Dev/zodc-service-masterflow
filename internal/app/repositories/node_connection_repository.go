package repositories

import (
	"context"
	"database/sql"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/table"
)

type NodeConnectionRepository struct{}

func NewNodeConnectionRepository() *NodeConnectionRepository {
	return &NodeConnectionRepository{}
}

func (r *NodeConnectionRepository) Create(ctx context.Context, tx *sql.Tx, connections []model.NodeConnections) error {
	NodeConnections := table.NodeConnections

	nodeConnectionInsertColumns := NodeConnections.AllColumns.Except(NodeConnections.CreatedAt, NodeConnections.UpdatedAt, NodeConnections.DeletedAt)

	nodeConnectionStmt := NodeConnections.INSERT(nodeConnectionInsertColumns).MODELS(connections)

	if err := nodeConnectionStmt.QueryContext(ctx, tx, &connections); err != nil {
		return err
	}

	return nil
}
