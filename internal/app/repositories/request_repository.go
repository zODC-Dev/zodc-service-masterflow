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

	rRequestStatement := postgres.SELECT(
		Requests.AllColumns,
	).FROM(
		Requests.
			LEFT_JOIN(WorkflowVerions, Requests.WorkflowVersionID.EQ(WorkflowVerions.ID)).
			LEFT_JOIN(Workflows, WorkflowVerions.WorkflowID.EQ(Workflows.ID)),
	).LIMIT(int64(requestQueryParam.PageSize)).OFFSET(int64(requestQueryParam.Page - 1))

	rRequestsConditions := []postgres.BoolExpression{
		Requests.UserID.EQ(postgres.Int32(userId)),
		Requests.IsTemplate.EQ(postgres.Bool(false)),
	}

	if requestQueryParam.WorkflowType != "" {
		rRequestsConditions = append(rRequestsConditions, Workflows.Type.EQ(postgres.String(requestQueryParam.WorkflowType)))
	}

	if requestQueryParam.Search != "" {
		rRequestsConditions = append(rRequestsConditions, postgres.LOWER(Requests.Title).LIKE(postgres.LOWER(postgres.String("%"+requestQueryParam.Search+"%"))))
	}

	if requestQueryParam.ProjectKey != "" {
		rRequestsConditions = append(rRequestsConditions, Workflows.ProjectKey.EQ(postgres.String(requestQueryParam.ProjectKey)))
	}

	if requestQueryParam.Status != "" {
		if requestQueryParam.Status == "ALL" {
			rRequestsConditions = append(rRequestsConditions, Requests.UserID.EQ(postgres.Int32(userId)))
		} else {
			rRequestsConditions = append(rRequestsConditions, Requests.Status.EQ(postgres.String(requestQueryParam.Status)))
		}
	}

	if requestQueryParam.SprintID != "" {
		sprintIdInt, err := strconv.Atoi(requestQueryParam.SprintID)
		if err != nil {
			return results.Count{}, []results.Request{}, err
		}
		rRequestsConditions = append(rRequestsConditions, Requests.SprintID.EQ(postgres.Int64(int64(sprintIdInt))))
	}

	var rRequests postgres.SelectTable
	if len(rRequestsConditions) > 0 {
		rRequests = rRequestStatement.WHERE(postgres.AND(rRequestsConditions...)).AsTable("rRequests")
	} else {
		rRequests = rRequestStatement.AsTable("rRequests")
	}

	requestId := Requests.ID.From(rRequests)
	requestWorkflowVersionId := Requests.WorkflowVersionID.From(rRequests)

	statement := postgres.SELECT(
		rRequests.AllColumns(),
		Workflows.AllColumns,
		WorkflowVerions.AllColumns,
		Nodes.AllColumns,
	).FROM(
		rRequests.
			LEFT_JOIN(Nodes, requestId.EQ(Nodes.RequestID)).
			LEFT_JOIN(WorkflowVerions, requestWorkflowVersionId.EQ(WorkflowVerions.ID)).
			LEFT_JOIN(Workflows, Workflows.ID.EQ(WorkflowVerions.WorkflowID)),
	).WHERE(
		Nodes.IsCurrent.EQ(postgres.Bool(true)),
	)

	result := []results.Request{}
	err := statement.QueryContext(ctx, db, &result)

	if err != nil {
		return results.Count{}, result, err
	}

	// Count
	statementCount := postgres.SELECT(
		Requests.ID,
	).FROM(
		Requests.
			LEFT_JOIN(WorkflowVerions, WorkflowVerions.ID.EQ(Requests.WorkflowVersionID)).
			LEFT_JOIN(Workflows, Workflows.ID.EQ(WorkflowVerions.WorkflowID)),
	)

	conditionsCount := []postgres.BoolExpression{
		Requests.UserID.EQ(postgres.Int32(userId)),
		Requests.IsTemplate.EQ(postgres.Bool(false)),
	}

	if requestQueryParam.WorkflowType != "" {
		conditionsCount = append(conditionsCount, Workflows.Type.EQ(postgres.String(requestQueryParam.WorkflowType)))
	}

	if requestQueryParam.ProjectKey != "" {
		conditionsCount = append(conditionsCount, Workflows.ProjectKey.EQ(postgres.String(requestQueryParam.ProjectKey)))
	}

	if requestQueryParam.Status != "" {
		if requestQueryParam.Status == "ALL" {
			conditionsCount = append(conditionsCount, Requests.UserID.EQ(postgres.Int32(userId)))
		} else {
			conditionsCount = append(conditionsCount, Requests.Status.EQ(postgres.String(requestQueryParam.Status)))
		}
	}

	if requestQueryParam.SprintID != "" {
		sprintIdInt, err := strconv.Atoi(requestQueryParam.SprintID)
		if err != nil {
			return results.Count{}, []results.Request{}, err
		}
		conditionsCount = append(conditionsCount, Requests.SprintID.EQ(postgres.Int64(int64(sprintIdInt))))
	}

	if len(conditionsCount) > 0 {
		statementCount = statementCount.WHERE(postgres.AND(conditionsCount...))
	}

	resultCount := []model.Requests{}
	err = statementCount.QueryContext(ctx, db, &resultCount)

	if err != nil {
		return results.Count{}, result, err
	}

	return results.Count{Count: int64(len(resultCount))}, result, err
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
	NodeForms := table.NodeForms
	FormData := table.FormData
	FormFieldData := table.FormFieldData
	FormTemplateFields := table.FormTemplateFields
	NodeFormApproveOrRejectUsers := table.NodeFormApproveOrRejectUsers
	Connections := table.Connections
	Categories := table.Categories

	statement := postgres.SELECT(
		Requests.AllColumns,
		WorkflowVersions.AllColumns,
		Workflows.AllColumns,
		Nodes.AllColumns,
		NodeForms.AllColumns,
		NodeFormApproveOrRejectUsers.AllColumns,
		Connections.AllColumns,
		Categories.AllColumns,
		FormData.AllColumns,
		FormFieldData.AllColumns,
		FormTemplateFields.AllColumns,
	).FROM(
		Requests.
			LEFT_JOIN(WorkflowVersions, WorkflowVersions.ID.EQ(Requests.WorkflowVersionID)).
			LEFT_JOIN(Workflows, Workflows.ID.EQ(WorkflowVersions.WorkflowID)).
			LEFT_JOIN(Nodes, Nodes.RequestID.EQ(Requests.ID)).
			LEFT_JOIN(NodeForms, NodeForms.NodeID.EQ(Nodes.ID)).
			LEFT_JOIN(NodeFormApproveOrRejectUsers, NodeFormApproveOrRejectUsers.NodeFormID.EQ(NodeForms.ID)).
			LEFT_JOIN(FormData, FormData.ID.EQ(Nodes.FormDataID)).
			LEFT_JOIN(FormFieldData, FormFieldData.FormDataID.EQ(FormData.ID)).
			LEFT_JOIN(FormTemplateFields, FormTemplateFields.ID.EQ(FormFieldData.FormTemplateFieldID)).
			LEFT_JOIN(Connections, Connections.RequestID.EQ(Requests.ID)).
			LEFT_JOIN(Categories, Workflows.CategoryID.EQ(Categories.ID)),
	).WHERE(
		Requests.ID.EQ(postgres.Int32(requestId)),
	).ORDER_BY(
		Nodes.Key.DESC(),
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

func (r *RequestRepository) CountRequestTaskByStatusAndUserIdAndQueryParams(ctx context.Context, db *sql.DB, userId int32, status constants.NodeStatus, queryparams queryparams.RequestTaskCount) (int, error) {
	Requests := table.Requests
	Workflows := table.Workflows
	WorkflowVersions := table.WorkflowVersions
	Nodes := table.Nodes

	statement := postgres.SELECT(
		Nodes.ID,
	).FROM(
		Nodes.
			LEFT_JOIN(
				Requests, Nodes.RequestID.EQ(Requests.ID),
			).
			LEFT_JOIN(
				WorkflowVersions, Requests.WorkflowVersionID.EQ(WorkflowVersions.ID),
			).
			LEFT_JOIN(
				Workflows, Workflows.ID.EQ(WorkflowVersions.WorkflowID),
			),
	)

	conditions := []postgres.BoolExpression{
		Nodes.AssigneeID.EQ(postgres.Int32(userId)),
		Requests.IsTemplate.EQ(postgres.Bool(false)),
		Nodes.Type.NOT_EQ(postgres.String(string(constants.NodeTypeStart))),
		Nodes.Type.NOT_EQ(postgres.String(string(constants.NodeTypeEnd))),
		Nodes.Type.NOT_EQ(postgres.String(string(constants.NodeTypeStory))),
	}

	if queryparams.ProjectKey != "" {
		conditions = append(conditions, Workflows.ProjectKey.EQ(postgres.String(queryparams.ProjectKey)))
	}

	if queryparams.Type != "" {
		conditions = append(conditions, Nodes.Type.EQ(postgres.String(queryparams.Type)))
	}

	if queryparams.WorkflowType != "" {
		conditions = append(conditions, Workflows.Type.EQ(postgres.String(queryparams.WorkflowType)))
	}

	if status != "" {
		conditions = append(conditions, Nodes.Status.EQ(postgres.String(string(status))))
	}

	if len(conditions) > 0 {
		statement = statement.WHERE(postgres.AND(conditions...))
	}

	result := []model.Nodes{}
	err := statement.QueryContext(ctx, db, &result)
	fmt.Println(statement.DebugSql())

	return len(result), err
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

func (r *RequestRepository) FindAllTasksByProject(ctx context.Context, db *sql.DB, userId int32, queryparams queryparams.RequestTaskProjectQueryParam) (int, []results.NodeResult, error) {
	Requests := table.Requests
	Workflows := table.Workflows
	WorkflowVersions := table.WorkflowVersions
	Nodes := table.Nodes

	statement := postgres.SELECT(
		Nodes.AllColumns,
		Requests.AllColumns,
	).FROM(
		Nodes.
			INNER_JOIN(
				Requests, Nodes.RequestID.EQ(Requests.ID),
			).
			INNER_JOIN(
				WorkflowVersions, Requests.WorkflowVersionID.EQ(WorkflowVersions.ID),
			).
			INNER_JOIN(
				Workflows, Workflows.ID.EQ(WorkflowVersions.WorkflowID),
			),
	).LIMIT(int64(queryparams.PageSize)).OFFSET(int64(queryparams.Page - 1))

	conditions := []postgres.BoolExpression{
		Nodes.AssigneeID.EQ(postgres.Int32(userId)),
		Requests.IsTemplate.EQ(postgres.Bool(false)),
		Nodes.Type.NOT_EQ(postgres.String(string(constants.NodeTypeStart))),
		Nodes.Type.NOT_EQ(postgres.String(string(constants.NodeTypeEnd))),
	}

	if queryparams.ProjectKey != "" {
		conditions = append(conditions, Workflows.ProjectKey.EQ(postgres.String(queryparams.ProjectKey)))
	}

	if queryparams.Status != "" {
		if queryparams.Status == "TODAY" {
			conditions = append(conditions, Nodes.IsCurrent.EQ(postgres.Bool(true)))
			conditions = append(conditions, Nodes.Status.NOT_EQ(postgres.String(string(constants.NodeStatusCompleted))))
		} else if queryparams.Status == "INCOMING" {
			conditions = append(conditions, Nodes.IsCurrent.EQ(postgres.Bool(false)))
			conditions = append(conditions, Nodes.Status.NOT_EQ(postgres.String(string(constants.NodeStatusCompleted))))
		} else {
			conditions = append(conditions, Nodes.Status.EQ(postgres.String(queryparams.Status)))
		}
	}

	if queryparams.Type != "" {
		conditions = append(conditions, Nodes.Type.EQ(postgres.String(queryparams.Type)))
	}

	if queryparams.WorkflowType != "" {
		conditions = append(conditions, Workflows.Type.EQ(postgres.String(queryparams.WorkflowType)))
	}

	if len(conditions) > 0 {
		statement = statement.WHERE(postgres.AND(conditions...))
	}

	result := []results.NodeResult{}
	err := statement.QueryContext(ctx, db, &result)
	fmt.Println(statement.DebugSql())
	if err != nil {
		return 0, result, err
	}

	//count
	statementCount := postgres.SELECT(
		Nodes.AllColumns,
		Requests.AllColumns,
	).FROM(
		Nodes.
			INNER_JOIN(
				Requests, Nodes.RequestID.EQ(Requests.ID),
			).
			INNER_JOIN(
				WorkflowVersions, Requests.WorkflowVersionID.EQ(WorkflowVersions.ID),
			).
			INNER_JOIN(
				Workflows, Workflows.ID.EQ(WorkflowVersions.WorkflowID),
			),
	).LIMIT(int64(queryparams.PageSize)).OFFSET(int64(queryparams.Page - 1))

	conditionsCount := []postgres.BoolExpression{
		Nodes.AssigneeID.EQ(postgres.Int32(userId)),
		Requests.IsTemplate.EQ(postgres.Bool(false)),
		Nodes.Type.NOT_EQ(postgres.String(string(constants.NodeTypeStart))),
		Nodes.Type.NOT_EQ(postgres.String(string(constants.NodeTypeEnd))),
	}

	if queryparams.ProjectKey != "" {
		conditionsCount = append(conditionsCount, Workflows.ProjectKey.EQ(postgres.String(queryparams.ProjectKey)))
	}

	if queryparams.Status != "" {
		if queryparams.Status == "TODAY" {
			conditionsCount = append(conditionsCount, Nodes.IsCurrent.EQ(postgres.Bool(true)))
			conditionsCount = append(conditionsCount, Nodes.Status.NOT_EQ(postgres.String(string(constants.NodeStatusCompleted))))
		} else if queryparams.Status == "INCOMING" {
			conditionsCount = append(conditionsCount, Nodes.IsCurrent.EQ(postgres.Bool(false)))
		} else {
			conditionsCount = append(conditionsCount, Requests.Status.EQ(postgres.String(queryparams.Status)))
		}
	}

	if queryparams.Type != "" {
		conditionsCount = append(conditionsCount, Nodes.Type.EQ(postgres.String(queryparams.Type)))
	}

	if queryparams.WorkflowType != "" {
		conditionsCount = append(conditionsCount, Workflows.Type.EQ(postgres.String(queryparams.WorkflowType)))
	}

	if len(conditions) > 0 {
		statementCount = statementCount.WHERE(postgres.AND(conditionsCount...))
	}

	count := []model.Nodes{}
	err = statementCount.QueryContext(ctx, db, &count)

	return len(count), result, err
}

func (r *RequestRepository) FindAllSubRequestByParentId(ctx context.Context, db *sql.DB, parentId int32, queryparams queryparams.RequestSubRequestQueryParam) (int, []results.RequestSubRequest, error) {
	Requests := table.Requests
	Nodes := table.Nodes
	Workflows := table.Workflows
	WorkflowVersions := table.WorkflowVersions

	var err error

	statement := Requests.SELECT(
		Requests.AllColumns,
		Nodes.AllColumns,
		Workflows.AllColumns,
	).FROM(
		Requests.
			LEFT_JOIN(Nodes, Requests.ID.EQ(Nodes.SubRequestID)).
			LEFT_JOIN(WorkflowVersions, WorkflowVersions.ID.EQ(Requests.WorkflowVersionID)).
			LEFT_JOIN(Workflows, WorkflowVersions.WorkflowID.EQ(Workflows.ID)),
	).WHERE(
		Requests.ParentID.EQ(postgres.Int32(parentId)).AND(Nodes.Type.IN(postgres.String(string(constants.NodeTypeTask)), postgres.String(string(constants.NodeTypeStory)))),
	).LIMIT(int64(queryparams.PageSize)).OFFSET(int64(queryparams.Page - 1))

	requests := []results.RequestSubRequest{}
	err = statement.QueryContext(ctx, db, &requests)
	if err != nil {
		return 0, requests, err
	}

	statementCount := Requests.SELECT(
		Requests.ID,
	).FROM(
		Requests,
	).WHERE(
		Requests.ParentID.EQ(postgres.Int32(parentId)),
	)

	count := []model.Requests{}
	err = statementCount.QueryContext(ctx, db, &count)

	return len(count), requests, err
}

func (r *RequestRepository) FindAllSubRequestByParentIdWithoutPagination(ctx context.Context, db *sql.DB, parentId int32) ([]results.RequestSubRequest, error) {
	Requests := table.Requests
	Workflows := table.Workflows
	WorkflowVersions := table.WorkflowVersions

	statement := Requests.SELECT(
		Requests.AllColumns,
		Workflows.AllColumns,
	).FROM(
		Requests.
			LEFT_JOIN(WorkflowVersions, WorkflowVersions.ID.EQ(Requests.WorkflowVersionID)).
			LEFT_JOIN(Workflows, WorkflowVersions.WorkflowID.EQ(Workflows.ID)),
	).WHERE(
		Requests.ParentID.EQ(postgres.Int32(parentId)),
	)

	requests := []results.RequestSubRequest{}
	err := statement.QueryContext(ctx, db, &requests)

	return requests, err
}

func (r *RequestRepository) RemoveNodesConnectionsStoriesByRequestId(ctx context.Context, tx *sql.Tx, requestId int32) error {
	Nodes := table.Nodes
	Connections := table.Connections

	statementNodes := Nodes.DELETE().WHERE(
		Nodes.RequestID.EQ(postgres.Int32(requestId)),
	)

	statementConnections := Connections.DELETE().WHERE(
		Connections.RequestID.EQ(postgres.Int32(requestId)),
	)

	var err error

	_, err = statementConnections.ExecContext(ctx, tx)
	if err != nil {
		return err
	}

	_, err = statementNodes.ExecContext(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}

func (r *RequestRepository) FindAllRequestCompletedFormByRequestId(ctx context.Context, db *sql.DB, requestId int32, page int, pageSize int) (int, []results.NodeFormCompletedResult, error) {
	NodeForms := table.NodeForms
	Nodes := table.Nodes
	NodeFormApproveOrRejectUsers := table.NodeFormApproveOrRejectUsers
	FormData := table.FormData
	FormFieldData := table.FormFieldData
	FormTemplateVersions := table.FormTemplateVersions
	FormTemplates := table.FormTemplates

	statement := NodeForms.SELECT(
		NodeForms.AllColumns,
		Nodes.AllColumns,
		NodeFormApproveOrRejectUsers.AllColumns,
		FormData.AllColumns,
		FormFieldData.AllColumns,
		FormTemplateVersions.AllColumns,
		FormTemplates.AllColumns,
	).FROM(
		NodeForms.
			LEFT_JOIN(Nodes, NodeForms.NodeID.EQ(Nodes.ID)).
			LEFT_JOIN(
				NodeFormApproveOrRejectUsers, NodeForms.ID.EQ(NodeFormApproveOrRejectUsers.NodeFormID),
			).
			LEFT_JOIN(FormData, NodeForms.DataID.EQ(FormData.ID)).
			LEFT_JOIN(FormFieldData, FormData.ID.EQ(FormFieldData.FormDataID)).
			//
			LEFT_JOIN(FormTemplateVersions, FormData.FormTemplateVersionID.EQ(FormTemplateVersions.ID)).
			LEFT_JOIN(FormTemplates, FormTemplateVersions.FormTemplateID.EQ(FormTemplates.ID)),
	).WHERE(
		Nodes.RequestID.EQ(postgres.Int32(requestId)).
			AND(NodeForms.IsSubmitted.EQ(postgres.Bool(true))),
	).LIMIT(int64(pageSize)).OFFSET(int64(page - 1))

	result := []results.NodeFormCompletedResult{}
	err := statement.QueryContext(ctx, db, &result)
	if err != nil {
		return 0, result, err
	}

	statementCount := NodeForms.SELECT(
		NodeForms.ID,
	).FROM(
		NodeForms.
			LEFT_JOIN(Nodes, NodeForms.NodeID.EQ(Nodes.ID)),
	).WHERE(
		Nodes.RequestID.EQ(postgres.Int32(requestId)).
			AND(NodeForms.IsSubmitted.EQ(postgres.Bool(true))),
	)

	count := []model.NodeForms{}
	err = statementCount.QueryContext(ctx, db, &count)

	return len(count), result, err
}

func (r *RequestRepository) FindAllRequestFileManagerByRequestId(ctx context.Context, db *sql.DB, requestId int32, page int, pageSize int) (int, []results.NodeFormResult, error) {
	Requests := table.Requests
	Nodes := table.Nodes
	NodeForms := table.NodeForms

	FormData := table.FormData
	FormFieldData := table.FormFieldData

	FormTemplateFields := table.FormTemplateFields

	statement := Requests.SELECT(
		NodeForms.AllColumns,
		FormFieldData.AllColumns,
	).FROM(
		Requests.
			LEFT_JOIN(Nodes, Requests.ID.EQ(Nodes.RequestID)).
			LEFT_JOIN(NodeForms, Nodes.ID.EQ(NodeForms.NodeID)).

			//
			LEFT_JOIN(FormData, FormData.ID.EQ(NodeForms.DataID)).
			LEFT_JOIN(FormFieldData, FormFieldData.FormDataID.EQ(FormData.ID)).
			LEFT_JOIN(FormTemplateFields, FormTemplateFields.ID.EQ(FormFieldData.FormTemplateFieldID)),
	).WHERE(
		Requests.ID.EQ(postgres.Int32(requestId)).
			AND(Nodes.Type.EQ(postgres.String(string(constants.NodeTypeInput)))).
			AND(NodeForms.IsSubmitted.EQ(postgres.Bool(true))).
			AND(FormTemplateFields.FieldType.EQ(postgres.String(string(constants.FormTemplateFieldTypeAttachment)))),
	).LIMIT(int64(pageSize)).OFFSET(int64(page - 1))

	result := []results.NodeFormResult{}
	err := statement.QueryContext(ctx, db, &result)

	return len(result), result, err

}
