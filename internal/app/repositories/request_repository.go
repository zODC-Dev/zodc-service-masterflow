package repositories

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/table"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/constants"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/results"
)

type RequestRepository struct{}

func NewRequestRepository() *RequestRepository {
	return &RequestRepository{}
}

func (r *RequestRepository) FindAllRequest(ctx context.Context, db *sql.DB, requestQueryParam queryparams.RequestQueryParam, userId int32) (results.Count, []results.Request, error) {
	Nodes := table.Nodes
	Requests := table.Requests
	WorkflowVerions := table.WorkflowVersions
	Workflows := table.Workflows

	statement := postgres.SELECT(
		Requests.AllColumns,
		Workflows.AllColumns,
		WorkflowVerions.AllColumns,
		Nodes.AllColumns,
	).FROM(
		Requests.
			LEFT_JOIN(Nodes, Requests.ID.EQ(Nodes.RequestID)).
			LEFT_JOIN(WorkflowVerions, Requests.WorkflowVersionID.EQ(WorkflowVerions.ID)).
			LEFT_JOIN(Workflows, Workflows.ID.EQ(WorkflowVerions.WorkflowID)),
	).WHERE(
		Requests.UserID.EQ(postgres.Int32(userId)),
	).LIMIT(int64(requestQueryParam.PageSize)).OFFSET(int64(requestQueryParam.Page))

	conditions := []postgres.BoolExpression{}

	// Filter out template requests
	conditions = append(conditions, Requests.IsTemplate.EQ(postgres.Bool(false)))

	if requestQueryParam.Search != "" {
		conditions = append(conditions, postgres.LOWER(Requests.Title).LIKE(postgres.LOWER(postgres.String("%"+requestQueryParam.Search+"%"))))
	}

	if requestQueryParam.ProjectKey != "" {
		conditions = append(conditions, Workflows.ProjectKey.EQ(postgres.String(requestQueryParam.ProjectKey)))
	}

	if requestQueryParam.Status != "" {
		if requestQueryParam.Status == "ALL" {
			statement = statement.WHERE(Nodes.AssigneeID.EQ(postgres.Int32(userId)))
		} else {
			conditions = append(conditions, Requests.Status.EQ(postgres.String(requestQueryParam.Status)))
		}
	}

	if requestQueryParam.SprintID != "" {
		sprintIdInt, err := strconv.Atoi(requestQueryParam.SprintID)
		if err != nil {
			return results.Count{}, []results.Request{}, err
		}
		conditions = append(conditions, Requests.SprintID.EQ(postgres.Int64(int64(sprintIdInt))))
	}

	if len(conditions) > 0 {
		statement = statement.WHERE(postgres.AND(conditions...))
	}

	result := []results.Request{}
	err := statement.QueryContext(ctx, db, &result)

	if err != nil {
		return results.Count{}, result, err
	}

	// Count
	statementCount := postgres.SELECT(
		postgres.COUNT(Requests.ID).AS("count"),
	).FROM(
		Requests,
	).WHERE(
		Requests.UserID.EQ(postgres.Int32(userId)),
	)

	count := results.Count{}

	err = statementCount.QueryContext(ctx, db, &count)
	if err != nil {
		return results.Count{}, result, err
	}

	return count, result, err
}

func (r *RequestRepository) FindRequestByNodeId(ctx context.Context, db *sql.DB, nodeId string) (model.Requests, error) {
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

func (r *RequestRepository) FindOneRequestByRequestId(ctx context.Context, db *sql.DB, requestId int32) (results.RequestDetail, error) {
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

func (r *RequestRepository) FindOneRequestByRequestIdWithNodeAssigneeId(ctx context.Context, db *sql.DB, requestId int32, assigneeId int32) (results.RequestDetail, error) {
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
		Requests.ID.EQ(postgres.Int32(requestId)).AND(Nodes.AssigneeID.EQ(postgres.Int32(assigneeId))),
	)

	result := results.RequestDetail{}

	err := statement.QueryContext(ctx, db, &result)

	return result, err
}

func (r *RequestRepository) FindOneRequestByRequestIdTx(ctx context.Context, tx *sql.Tx, requestId int32) (results.RequestDetail, error) {
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

	err := statement.QueryContext(ctx, tx, &result)

	return result, err
}

func (r *RequestRepository) UpdateRequest(ctx context.Context, tx *sql.Tx, request model.Requests) error {
	Requests := table.Requests

	request.UpdatedAt = time.Now()

	columns := Requests.AllColumns.Except(Requests.ID, Requests.CreatedAt, Requests.DeletedAt)

	statement := Requests.UPDATE(columns).MODEL(request).WHERE(Requests.ID.EQ(postgres.Int32(request.ID))).RETURNING(Requests.ID)

	err := statement.QueryContext(ctx, tx, &request)

	return err

}

func (r *RequestRepository) CreateRequest(ctx context.Context, tx *sql.Tx, request model.Requests) (model.Requests, error) {
	Requests := table.Requests

	columns := Requests.AllColumns.Except(Requests.ID, Requests.CreatedAt, Requests.UpdatedAt, Requests.DeletedAt, Requests.Key)

	statement := Requests.INSERT(columns).MODEL(request).RETURNING(Requests.ID)

	err := statement.QueryContext(ctx, tx, &request)

	return request, err
}

func (r *RequestRepository) CountRequestByStatusAndUserId(ctx context.Context, db *sql.DB, userId int32, status constants.RequestStatus) (int64, error) {
	Requests := table.Requests

	statement := postgres.SELECT(
		postgres.COUNT(Requests.ID).AS("count"),
	).FROM(
		Requests,
	).WHERE(
		Requests.UserID.EQ(postgres.Int32(userId)),
	)

	conditions := []postgres.BoolExpression{}

	if status != "" {
		conditions = append(conditions, Requests.Status.EQ(postgres.String(string(status))))
	}

	if len(conditions) > 0 {
		statement = statement.WHERE(postgres.AND(conditions...))
	}

	count := results.Count{}
	err := statement.QueryContext(ctx, db, &count)

	return count.Count, err

}

func (r *RequestRepository) CountUserAppendInRequestAndNodeUserId(ctx context.Context, db *sql.DB, userId int32) (int64, error) {
	Requests := table.Requests
	Nodes := table.Nodes

	statement := postgres.SELECT(
		postgres.COUNT(Requests.ID).AS("count"),
	).FROM(
		Requests.
			LEFT_JOIN(Nodes, Nodes.RequestID.EQ(Requests.ID)),
	).WHERE(
		Requests.UserID.EQ(postgres.Int32(userId)).
			OR(Nodes.AssigneeID.EQ(postgres.Int32(userId))),
	)

	count := results.Count{}
	err := statement.QueryContext(ctx, db, &count)

	return count.Count, err

}

func (r *RequestRepository) FindAllChildrenRequestByRequestId(ctx context.Context, db *sql.DB, requestId int32) ([]model.Requests, error) {
	Requests := table.Requests

	statement := postgres.SELECT(
		Requests.AllColumns,
	).FROM(
		Requests,
	).WHERE(
		Requests.ParentID.EQ(postgres.Int32(requestId)),
	)

	result := []model.Requests{}
	err := statement.QueryContext(ctx, db, &result)

	return result, err
}
