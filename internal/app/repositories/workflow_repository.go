package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/table"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/results"
)

type WorkflowRepository struct{}

func NewWorkflowRepository() *WorkflowRepository {
	return &WorkflowRepository{}
}

func (r *WorkflowRepository) CreateWorkflow(ctx context.Context, tx *sql.Tx, workflow model.Workflows) (model.Workflows, error) {
	Workflows := table.Workflows

	columns := Workflows.AllColumns.Except(Workflows.ID, Workflows.CreatedAt, Workflows.UpdatedAt, Workflows.DeletedAt)

	statement := Workflows.INSERT(columns).MODEL(workflow).RETURNING(Workflows.AllColumns)

	err := statement.QueryContext(ctx, tx, &workflow)

	return workflow, err
}

// Create

func (r *WorkflowRepository) CreateWorkflowVersion(ctx context.Context, tx *sql.Tx, workflowVersion model.WorkflowVersions) (model.WorkflowVersions, error) {
	WorkflowVersions := table.WorkflowVersions

	columns := WorkflowVersions.AllColumns.Except(WorkflowVersions.ID, WorkflowVersions.CreatedAt, WorkflowVersions.UpdatedAt, WorkflowVersions.DeletedAt)

	statement := WorkflowVersions.INSERT(columns).MODEL(workflowVersion).RETURNING(WorkflowVersions.ID)

	err := statement.QueryContext(ctx, tx, &workflowVersion)

	return workflowVersion, err
}

func (r *WorkflowRepository) CreateRequest(ctx context.Context, tx *sql.Tx, request model.Requests) (model.Requests, error) {
	Requests := table.Requests

	columns := Requests.AllColumns.Except(Requests.ID, Requests.CreatedAt, Requests.UpdatedAt, Requests.DeletedAt, Requests.Key)

	statement := Requests.INSERT(columns).MODEL(request).RETURNING(Requests.ID)

	err := statement.QueryContext(ctx, tx, &request)

	return request, err
}

func (r *WorkflowRepository) CreateWorkflowNodes(ctx context.Context, tx *sql.Tx, nodes []model.Nodes) error {
	Nodes := table.Nodes

	columns := Nodes.AllColumns.Except(Nodes.CreatedAt, Nodes.UpdatedAt, Nodes.DeletedAt)

	statement := Nodes.INSERT(columns).MODELS(nodes)

	err := statement.QueryContext(ctx, tx, &nodes)

	return err
}

func (r *WorkflowRepository) CreateWorkflowConnections(ctx context.Context, tx *sql.Tx, connections []model.Connections) error {
	Connections := table.Connections

	columns := Connections.AllColumns.Except(Connections.CreatedAt, Connections.UpdatedAt, Connections.DeletedAt)

	statement := Connections.INSERT(columns).MODELS(connections)

	err := statement.QueryContext(ctx, tx, &connections)

	if err != nil {
		fmt.Println(statement.DebugSql())
	}

	return err
}

// Find

func (r *WorkflowRepository) FindAllWorkflowTemplates(ctx context.Context, db *sql.DB, workflowTemplateQueryParams queryparams.WorkflowQueryParam) ([]results.WorkflowTemplate, error) {
	Workflows := table.Workflows
	WorkflowVersions := table.WorkflowVersions
	Requests := table.Requests
	Categories := table.Categories

	statement := postgres.SELECT(
		Workflows.AllColumns,
		WorkflowVersions.AllColumns,
		Categories.AllColumns,
		Requests.AllColumns,
	).FROM(
		Workflows.
			LEFT_JOIN(WorkflowVersions, WorkflowVersions.WorkflowID.EQ(Workflows.ID)).
			LEFT_JOIN(Categories, Workflows.CategoryID.EQ(Categories.ID)).
			LEFT_JOIN(Requests, Requests.WorkflowVersionID.EQ(WorkflowVersions.ID)),
	).WHERE(
		WorkflowVersions.Version.EQ(Workflows.Currentversion),
	)

	conditions := []postgres.BoolExpression{}

	if workflowTemplateQueryParams.Search != "" {
		conditions = append(conditions, postgres.LOWER(Workflows.Title).LIKE(postgres.LOWER(postgres.String("%"+workflowTemplateQueryParams.Search+"%"))))
	}

	if workflowTemplateQueryParams.Type != "" {
		conditions = append(conditions, Workflows.Type.EQ(postgres.String(workflowTemplateQueryParams.Type)))
	}

	if workflowTemplateQueryParams.CategoryID != "" {
		categoryIdInt, err := strconv.Atoi(workflowTemplateQueryParams.CategoryID)
		if err != nil {
			return []results.WorkflowTemplate{}, err
		}

		conditions = append(conditions, Workflows.CategoryID.EQ(postgres.Int32(int32(categoryIdInt))))
	}

	if workflowTemplateQueryParams.ProjectKey != "" {
		conditions = append(conditions, Workflows.ProjectKey.EQ(postgres.String(workflowTemplateQueryParams.ProjectKey)))
	}

	if workflowTemplateQueryParams.HasSubWorkflow != "" {
		hasSubWorkflowBool, err := strconv.ParseBool(workflowTemplateQueryParams.HasSubWorkflow)
		if err != nil {
			return []results.WorkflowTemplate{}, err
		}

		conditions = append(conditions, WorkflowVersions.HasSubWorkflow.EQ(postgres.Bool(hasSubWorkflowBool)))
	}

	if workflowTemplateQueryParams.IsArchived != "" {
		hasIsArchivedBool, err := strconv.ParseBool(workflowTemplateQueryParams.IsArchived)
		if err != nil {
			return []results.WorkflowTemplate{}, err
		}

		conditions = append(conditions, Workflows.IsArchived.EQ(postgres.Bool(hasIsArchivedBool)))
	}

	if len(conditions) > 0 {
		statement = statement.WHERE(postgres.AND(conditions...))
	}

	result := []results.WorkflowTemplate{}

	err := statement.QueryContext(ctx, db, &result)

	return result, err
}

func (r *WorkflowRepository) FindAllConnectionByREquestId(ctx context.Context, db *sql.DB, requestId int32) ([]model.Connections, error) {
	Connections := table.Connections

	statement := postgres.SELECT(
		Connections.AllColumns,
	).FROM(
		Connections,
	).WHERE(
		Connections.RequestID.EQ(postgres.Int32(requestId)),
	)

	results := []model.Connections{}

	err := statement.QueryContext(ctx, db, &results)

	return results, err
}

func (r *WorkflowRepository) FindAllNodeByRequestId(ctx context.Context, db *sql.DB, requestId int32) ([]model.Nodes, error) {
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

func (r *WorkflowRepository) FindOneRequestByRequestId(ctx context.Context, db *sql.DB, requestId int32) (results.RequestDetail, error) {
	Workflows := table.Workflows
	WorkflowVersions := table.WorkflowVersions
	Requests := table.Requests
	Nodes := table.Nodes
	Connections := table.Connections
	Categories := table.Categories

	statement := postgres.SELECT(
		Requests.AllColumns,
		WorkflowVersions.AllColumns,
		Workflows.AllColumns,
		Nodes.AllColumns,
		Connections.AllColumns,
		Categories.AllColumns,
	).FROM(
		Requests.
			LEFT_JOIN(WorkflowVersions, WorkflowVersions.ID.EQ(Requests.WorkflowVersionID)).
			LEFT_JOIN(Workflows, Workflows.ID.EQ(WorkflowVersions.WorkflowID)).
			LEFT_JOIN(Nodes, Nodes.RequestID.EQ(Requests.ID)).
			LEFT_JOIN(Connections, Connections.RequestID.EQ(Requests.ID)).
			LEFT_JOIN(Categories, Workflows.CategoryID.EQ(Categories.ID)),
	).WHERE(
		Requests.ID.EQ(postgres.Int32(requestId)),
	)

	result := results.RequestDetail{}

	err := statement.QueryContext(ctx, db, &result)

	return result, err
}

func (r *WorkflowRepository) FindOneNodeByNodeId(ctx context.Context, db *sql.DB, nodeId string) (model.Nodes, error) {
	Nodes := table.Nodes

	statement := postgres.SELECT(Nodes.AllColumns).
		FROM(Nodes).
		WHERE(Nodes.ID.EQ(postgres.String(nodeId)))

	result := model.Nodes{}
	err := statement.QueryContext(ctx, db, result)

	return result, err
}

func (r *WorkflowRepository) FindConnectionsWithToNodesByFromNodeId(ctx context.Context, db *sql.DB, fromNodeId string) ([]results.ConnectionWithNode, error) {
	Connections := table.Connections
	Nodes := table.Nodes

	statement := postgres.SELECT(
		Connections.AllColumns,
		Nodes.AllColumns,
	).FROM(
		Connections.
			INNER_JOIN(
				Nodes, Connections.ToNodeID.EQ(Nodes.ID),
			),
	).WHERE(Connections.FromNodeID.EQ(postgres.String(fromNodeId)))

	result := []results.ConnectionWithNode{}
	err := statement.QueryContext(ctx, db, &result)

	return result, err
}

func (r *WorkflowRepository) FindConnectionsByToNodeId(ctx context.Context, db *sql.DB, toNodeId string) ([]model.Connections, error) {
	Connections := table.Connections

	statement := postgres.SELECT(Connections.AllColumns).FROM(Connections).WHERE(Connections.ToNodeID.EQ(postgres.String(toNodeId)))

	result := []model.Connections{}
	err := statement.QueryContext(ctx, db, &result)

	return result, err
}

func (r *WorkflowRepository) FindRequestByNodeId(ctx context.Context, db *sql.DB, nodeId string) (model.Requests, error) {
	Nodes := table.Nodes
	Requests := table.Requests

	statement := postgres.SELECT(
		Requests.AllColumns,
	).FROM(
		Nodes.INNER_JOIN(
			Requests, Requests.ID.EQ(Nodes.RequestID),
		),
	).WHERE(Nodes.ID.EQ(postgres.String(nodeId)))

	result := model.Requests{}
	err := statement.QueryContext(ctx, db, &result)

	return result, err
}

func (r *WorkflowRepository) FindAllRequest(ctx context.Context, db *sql.DB, requestQueryParam queryparams.RequestQueryParam) ([]results.Request, error) {
	Nodes := table.Nodes
	Requests := table.Requests

	statement := postgres.SELECT(
		Requests.AllColumns,
	).FROM(
		Nodes.INNER_JOIN(
			Requests, Requests.ID.EQ(Nodes.RequestID),
		),
	).LIMIT(int64(requestQueryParam.PageSize)).OFFSET(int64(requestQueryParam.Page))

	conditions := []postgres.BoolExpression{}

	if requestQueryParam.Search != "" {
		conditions = append(conditions, postgres.LOWER(Requests.Title).LIKE(postgres.LOWER(postgres.String("%"+requestQueryParam.Search+"%"))))
	}

	if len(conditions) > 0 {
		statement = statement.WHERE(postgres.AND(conditions...))
	}

	result := []results.Request{}
	err := statement.QueryContext(ctx, db, &result)

	return result, err
}

func (r *WorkflowRepository) FindAllRequestCount(ctx context.Context, db *sql.DB) (int, error) {
	Requests := table.Requests

	statement := postgres.SELECT(
		postgres.COUNT(Requests.ID),
	).FROM(Requests)

	var count int
	err := statement.QueryContext(ctx, db, &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Update
func (r *WorkflowRepository) UpdateNode(ctx context.Context, tx *sql.Tx, node model.Nodes) error {
	Nodes := table.Nodes

	node.UpdatedAt = time.Now()

	columns := Nodes.AllColumns.Except(Nodes.ID, Nodes.CreatedAt, Nodes.DeletedAt)

	statement := Nodes.UPDATE(columns).MODEL(node).WHERE(Nodes.ID.EQ(postgres.String(node.ID)))

	err := statement.QueryContext(ctx, tx, &node)

	return err
}

func (r *WorkflowRepository) UpdateConnection(ctx context.Context, tx *sql.Tx, connection model.Connections) error {
	Connections := table.Connections

	connection.UpdatedAt = time.Now()

	columns := Connections.AllColumns.Except(Connections.ID, Connections.CreatedAt, Connections.DeletedAt)

	statment := Connections.UPDATE(columns).MODEL(connection)

	err := statment.QueryContext(ctx, tx, &connection)

	return err
}

func (r *WorkflowRepository) UpdateRequest(ctx context.Context, tx *sql.Tx, request model.Requests) error {
	Requests := table.Requests

	request.UpdatedAt = time.Now()

	columns := Requests.AllColumns.Except(Requests.ID, Requests.CreatedAt, Requests.DeletedAt)

	statement := Requests.UPDATE(columns).MODEL(request).WHERE(Requests.ID.EQ(postgres.Int32(request.ID)))

	err := statement.QueryContext(ctx, tx, &request)

	return err

}
