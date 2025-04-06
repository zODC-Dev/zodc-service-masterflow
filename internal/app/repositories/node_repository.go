package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
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

	columns := NodeConditionDestinations.AllColumns.Except(NodeConditionDestinations.ID, NodeConditionDestinations.CreatedAt, NodeConditionDestinations.UpdatedAt, NodeConditionDestinations.DeletedAt)

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

func (r *NodeRepository) CreateNodeForms(ctx context.Context, tx *sql.Tx, nodeForms []model.NodeForms) error {
	NodeForms := table.NodeForms

	columns := NodeForms.AllColumns.Except(NodeForms.ID, NodeForms.CreatedAt, NodeForms.UpdatedAt, NodeForms.DeletedAt)

	statement := NodeForms.INSERT(columns).MODELS(nodeForms)

	err := statement.QueryContext(ctx, tx, &nodeForms)

	return err
}

func (r *NodeRepository) UpdateNodePlannedTimes(ctx context.Context, tx *sql.Tx, nodeTimesUpdates []struct {
	NodeId           string
	PlannedStartTime time.Time
	PlannedEndTime   time.Time
}) error {
	if len(nodeTimesUpdates) == 0 {
		return nil
	}

	slog.Info("Starting to update planned times in database", "count", len(nodeTimesUpdates))

	// Thực hiện các updates riêng lẻ trong cùng một transaction
	for i, update := range nodeTimesUpdates {
		Nodes := table.Nodes

		// Log thông tin chi tiết của mỗi node đang được cập nhật
		slog.Info("Updating node planned times",
			"index", i,
			"nodeId", update.NodeId,
			"startTime", update.PlannedStartTime.Format(time.RFC3339),
			"endTime", update.PlannedEndTime.Format(time.RFC3339))

		// Tạo model để cập nhật
		nodeModel := model.Nodes{
			ID:               update.NodeId,
			PlannedStartTime: &update.PlannedStartTime,
			PlannedEndTime:   &update.PlannedEndTime,
			UpdatedAt:        time.Now(),
		}

		// Cập nhật theo cách sử dụng MODEL
		statement := Nodes.UPDATE(
			Nodes.PlannedStartTime,
			Nodes.PlannedEndTime,
			Nodes.UpdatedAt,
		).MODEL(nodeModel).WHERE(Nodes.ID.EQ(postgres.String(update.NodeId)))

		// Hiển thị SQL statement để debug
		sql, args := statement.Sql()
		slog.Info("SQL statement", "sql", sql, "args", args)

		result, err := statement.ExecContext(ctx, tx)
		if err != nil {
			slog.Error("Failed to update node", "nodeId", update.NodeId, "error", err)
			return fmt.Errorf("failed to update planned times for node %s: %w", update.NodeId, err)
		}

		// Kiểm tra số lượng row bị ảnh hưởng
		rowsAffected, _ := result.RowsAffected()
		slog.Info("Update result", "nodeId", update.NodeId, "rowsAffected", rowsAffected)

		if rowsAffected == 0 {
			slog.Warn("No rows affected when updating node", "nodeId", update.NodeId)
		}
	}

	slog.Info("Completed updating planned times in database")
	return nil
}
