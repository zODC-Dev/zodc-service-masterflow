package repositories

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/table"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/constants"
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

func (r *WorkflowRepository) CreateWorkflowVersion(ctx context.Context, tx *sql.Tx, workflowVersion model.WorkflowVersions) (model.WorkflowVersions, error) {
	WorkflowVersions := table.WorkflowVersions

	columns := WorkflowVersions.AllColumns.Except(WorkflowVersions.ID, WorkflowVersions.CreatedAt, WorkflowVersions.UpdatedAt, WorkflowVersions.DeletedAt, WorkflowVersions.IsArchived)

	statement := WorkflowVersions.INSERT(columns).MODEL(workflowVersion).RETURNING(WorkflowVersions.ID)

	err := statement.QueryContext(ctx, tx, &workflowVersion)

	return workflowVersion, err
}

func (r *WorkflowRepository) CreateWorkflowNodes(ctx context.Context, tx *sql.Tx, workflowNodes []model.WorkflowNodes) error {
	WorkflowNodes := table.WorkflowNodes

	columns := WorkflowNodes.AllColumns.Except(WorkflowNodes.CreatedAt, WorkflowNodes.UpdatedAt, WorkflowNodes.DeletedAt)

	statement := WorkflowNodes.INSERT(columns).MODELS(workflowNodes)

	err := statement.QueryContext(ctx, tx, &workflowNodes)

	return err
}

func (r *WorkflowRepository) CreateWorkflowConnections(ctx context.Context, tx *sql.Tx, workflowConnections []model.WorkflowConnections) error {
	WorkflowConnections := table.WorkflowConnections

	columns := WorkflowConnections.AllColumns.Except(WorkflowConnections.CreatedAt, WorkflowConnections.UpdatedAt, WorkflowConnections.DeletedAt)

	statement := WorkflowConnections.INSERT(columns).MODELS(workflowConnections)

	err := statement.QueryContext(ctx, tx, &workflowConnections)

	return err
}

func (r *WorkflowRepository) FindAllWorkflowTemplates(ctx context.Context, db *sql.DB) ([]results.WorkflowTemplateResult, error) {
	Workflows := table.Workflows
	WorkflowVersions := table.WorkflowVersions
	Categories := table.Categories

	statement := postgres.SELECT(
		Workflows.AllColumns,
		WorkflowVersions.AllColumns,
		Categories.AllColumns,
	).FROM(
		Workflows.
			INNER_JOIN(WorkflowVersions, WorkflowVersions.WorkflowID.EQ(Workflows.ID).
				AND(WorkflowVersions.IsArchived.EQ(postgres.Bool(false))),
			).
			INNER_JOIN(Categories, Workflows.CategoryID.EQ(Categories.ID)),
	)

	result := []results.WorkflowTemplateResult{}

	err := statement.QueryContext(ctx, db, &result)

	return result, err
}

func (r *WorkflowRepository) FindAllConnectionByWorkflowVersionId(ctx context.Context, db *sql.DB, workflowVersionId int32) ([]model.WorkflowConnections, error) {
	WorkflowConnections := table.WorkflowConnections

	statement := postgres.SELECT(
		WorkflowConnections.AllColumns,
	).FROM(
		WorkflowConnections,
	).WHERE(
		WorkflowConnections.WorkflowVersionID.EQ(postgres.Int32(workflowVersionId)),
	)

	results := []model.WorkflowConnections{}

	err := statement.QueryContext(ctx, db, &results)

	return results, err
}

func (r *WorkflowRepository) FindAllNodeByWorkflowVersionId(ctx context.Context, db *sql.DB, workflowVersionId int32) ([]model.WorkflowNodes, error) {
	WorkflowNodes := table.WorkflowNodes

	statement := postgres.SELECT(
		WorkflowNodes.AllColumns,
	).FROM(
		WorkflowNodes,
	).WHERE(
		WorkflowNodes.WorkflowVersionID.EQ(postgres.Int32(workflowVersionId)),
	)

	results := []model.WorkflowNodes{}

	err := statement.QueryContext(ctx, db, &results)

	return results, err
}

func (r *WorkflowRepository) FindOneWorkflowDetailByWorkflowVersionId(ctx context.Context, db *sql.DB, workflowVersionId int32) (results.WorkflowDetailResult, error) {
	Workflows := table.Workflows
	WorkflowVersions := table.WorkflowVersions
	WorkflowNodes := table.WorkflowNodes
	WorkflowConnections := table.WorkflowConnections

	statement := postgres.SELECT(
		Workflows.AllColumns,
		WorkflowVersions.AllColumns,
		WorkflowNodes.AllColumns,
		WorkflowConnections.AllColumns,
	).FROM(
		Workflows.
			INNER_JOIN(WorkflowVersions, WorkflowVersions.WorkflowID.EQ(Workflows.ID).
				AND(WorkflowVersions.IsArchived.EQ(postgres.Bool(false))).
				AND(WorkflowVersions.ID.EQ(postgres.Int32(workflowVersionId))),
			).
			INNER_JOIN(WorkflowNodes, WorkflowNodes.WorkflowVersionID.EQ(WorkflowVersions.ID)).
			INNER_JOIN(WorkflowConnections, WorkflowConnections.WorkflowVersionID.EQ(WorkflowVersions.ID)),
	).WHERE(
		Workflows.Type.EQ(postgres.String(string(constants.WorkflowGeneral))),
	)

	result := results.WorkflowDetailResult{}

	err := statement.QueryContext(ctx, db, &result)

	return result, err
}
