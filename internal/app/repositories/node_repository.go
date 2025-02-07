package repositories

import (
	"context"
	"database/sql"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/model"
	. "github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/table"
)

type NodeRepository struct{}

func NewNodeRepository() *NodeRepository {
	return &NodeRepository{}
}

func (r *NodeRepository) Create(ctx context.Context, tx *sql.Tx, nodes []model.Nodes) error {

	nodesInsertColumns := Nodes.AllColumns.Except(Nodes.CreatedAt, Nodes.UpdatedAt, Nodes.DeletedAt)

	nodesStmt := Nodes.INSERT(nodesInsertColumns).MODELS(nodes).RETURNING(Nodes.ID)

	nodeModels := []model.Nodes{}

	if err := nodesStmt.QueryContext(ctx, tx, &nodeModels); err != nil {
		return err
	}

	return nil
}
