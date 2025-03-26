package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/table"
)

type NodeRepository struct{}

func NewNodeRepository() *NodeRepository {
	return &NodeRepository{}
}

func (r *NodeRepository) UpdateNode(ctx context.Context, tx *sql.Tx, node model.Nodes) error {
	Nodes := table.Nodes

	node.UpdatedAt = time.Now()

	columns := Nodes.AllColumns.Except(Nodes.ID, Nodes.CreatedAt, Nodes.DeletedAt)

	statement := Nodes.UPDATE(columns).MODEL(node).WHERE(Nodes.ID.EQ(postgres.String(node.ID)))

	err := statement.QueryContext(ctx, tx, &node)

	return err
}

func (r *NodeRepository) FindOneNodeByNodeId(ctx context.Context, db *sql.DB, nodeId string) (model.Nodes, error) {
	Nodes := table.Nodes

	statement := postgres.SELECT(Nodes.AllColumns).
		FROM(Nodes).
		WHERE(Nodes.ID.EQ(postgres.String(nodeId)))

	result := model.Nodes{}
	err := statement.QueryContext(ctx, db, result)

	return result, err
}

func (r *NodeRepository) FindAllNodeByRequestId(ctx context.Context, db *sql.DB, requestId int32) ([]model.Nodes, error) {
	WorkflowNodes := table.Nodes

	statement := postgres.SELECT(
		WorkflowNodes.AllColumns,
	).FROM(
		WorkflowNodes,
	).WHERE(
		WorkflowNodes.RequestID.EQ(postgres.Int32(requestId)),
	)

	results := []model.Nodes{}

	err := statement.QueryContext(ctx, db, &results)

	return results, err
}

func (r *NodeRepository) CreateNodes(ctx context.Context, tx *sql.Tx, nodes []model.Nodes) error {
	Nodes := table.Nodes

	columns := Nodes.AllColumns.Except(Nodes.CreatedAt, Nodes.UpdatedAt, Nodes.DeletedAt)

	statement := Nodes.INSERT(columns).MODELS(nodes)

	err := statement.QueryContext(ctx, tx, &nodes)

	return err
}
