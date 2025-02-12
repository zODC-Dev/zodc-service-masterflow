package repositories

import (
	"context"
	"database/sql"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/table"
)

type NodeGroupRepository struct{}

func NewNodeGrouprepository() *NodeGroupRepository {
	return &NodeGroupRepository{}
}

func (r *NodeGroupRepository) Create(ctx context.Context, tx *sql.Tx, nodeGroups []model.NodeGroups) error {
	NodeGroups := table.NodeGroups

	nodeGroupInsertColumns := NodeGroups.AllColumns.Except(NodeGroups.DeletedAt, NodeGroups.CreatedAt, NodeGroups.UpdatedAt)

	stmt := NodeGroups.INSERT(nodeGroupInsertColumns).MODELS(nodeGroups)

	if err := stmt.QueryContext(ctx, tx, &nodeGroups); err != nil {
		return err
	}

	return nil
}
