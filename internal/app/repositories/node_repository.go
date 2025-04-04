package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/table"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/results"
)

type NodeRepository struct{}

func NewNodeRepository() *NodeRepository {
	return &NodeRepository{}
}

func (r *NodeRepository) UpdateNode(ctx context.Context, tx *sql.Tx, node model.Nodes) error {
	Nodes := table.Nodes

	node.UpdatedAt = time.Now()

	columns := Nodes.AllColumns.Except(Nodes.ID, Nodes.CreatedAt, Nodes.DeletedAt)

	statement := Nodes.UPDATE(columns).MODEL(node).WHERE(Nodes.ID.EQ(postgres.String(node.ID))).RETURNING(Nodes.ID)

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

func (r *NodeRepository) FindAllNodeByRequestIdWithPagination(ctx context.Context, db *sql.DB, requestId int32, requestTaskQueryParam queryparams.RequestTaskQueryParam, userId *int32) ([]model.Nodes, error) {
	Nodes := table.Nodes
	Requests := table.Requests

	statement := postgres.SELECT(
		Nodes.AllColumns,
		Requests.AllColumns,
	).FROM(
		Nodes,
	).WHERE(
		Nodes.RequestID.EQ(postgres.Int32(requestId)),
	).LIMIT(int64(requestTaskQueryParam.PageSize)).OFFSET(int64(requestTaskQueryParam.Page * requestTaskQueryParam.PageSize))

	if userId != nil {
		statement = statement.WHERE(
			Nodes.AssigneeID.EQ(postgres.Int32(*userId)),
		)
	}

	results := []model.Nodes{}

	err := statement.QueryContext(ctx, db, &results)

	return results, err
}

func (r *NodeRepository) CountAllNodeByRequestId(ctx context.Context, db *sql.DB, requestId int32) (results.Count, error) {
	Nodes := table.Nodes

	statement := postgres.SELECT(
		postgres.COUNT(Nodes.ID).AS("count"),
	).FROM(
		Nodes,
	).WHERE(
		Nodes.RequestID.EQ(postgres.Int32(requestId)),
	)

	count := results.Count{}
	err := statement.QueryContext(ctx, db, &count)

	return count, err
}

func (r *NodeRepository) CreateNodes(ctx context.Context, tx *sql.Tx, nodes []model.Nodes) error {
	Nodes := table.Nodes

	columns := Nodes.AllColumns.Except(Nodes.CreatedAt, Nodes.UpdatedAt, Nodes.DeletedAt, Nodes.Key)

	statement := Nodes.INSERT(columns).MODELS(nodes)

	err := statement.QueryContext(ctx, tx, &nodes)

	return err
}

func (r *NodeRepository) CreateNodeConditionDestinations(ctx context.Context, tx *sql.Tx, nodeConditionDestinations []model.NodeConditionDestinations) error {
	NodeConditionDestinations := table.NodeConditionDestinations

	columns := NodeConditionDestinations.AllColumns.Except(NodeConditionDestinations.CreatedAt, NodeConditionDestinations.UpdatedAt, NodeConditionDestinations.DeletedAt)

	statement := NodeConditionDestinations.INSERT(columns).MODELS(nodeConditionDestinations)

	err := statement.QueryContext(ctx, tx, &nodeConditionDestinations)

	return err
}

func (r *NodeRepository) UpdateJiraKey(ctx context.Context, tx *sql.Tx, nodeId string, jiraKey string) error {
	Nodes := table.Nodes

	statement := Nodes.UPDATE(Nodes.JiraKey).
		SET(postgres.String(jiraKey)).
		WHERE(Nodes.ID.EQ(postgres.String(nodeId)))

	_, err := statement.ExecContext(ctx, tx)
	return err
}
