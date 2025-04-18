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

// Find

func (r *WorkflowRepository) FindAllWorkflowTemplates(ctx context.Context, db *sql.DB, workflowTemplateQueryParams queryparams.WorkflowQueryParam, projects []string) ([]results.WorkflowTemplate, error) {
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
	)

	conditions := []postgres.BoolExpression{
		Requests.IsTemplate.EQ(postgres.Bool(true)),
		WorkflowVersions.Version.EQ(Workflows.CurrentVersion),
	}

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

	// Filter Product
	if workflowTemplateQueryParams.Type == string(constants.WorkflowTypeProject) && len(projects) > 0 {
		projectExpressions := make([]postgres.Expression, len(projects))
		for i, project := range projects {
			projectExpressions[i] = postgres.String(project)
		}

		conditions = append(conditions, Workflows.ProjectKey.IN(projectExpressions...))
	}

	if len(conditions) > 0 {
		statement = statement.WHERE(postgres.AND(conditions...))
	}

	result := []results.WorkflowTemplate{}

	err := statement.QueryContext(ctx, db, &result)

	return result, err
}

func (r *WorkflowRepository) FindOneWorkflowByWorkflowId(ctx context.Context, db *sql.DB, workflowId int32) (results.Workflow, error) {
	Workflows := table.Workflows
	WorkflowVersions := table.WorkflowVersions
	Requests := table.Requests

	statement := Workflows.SELECT(
		Workflows.AllColumns,
		WorkflowVersions.AllColumns,
		Requests.AllColumns,
	).FROM(
		Workflows.
			LEFT_JOIN(WorkflowVersions, WorkflowVersions.WorkflowID.EQ(Workflows.ID)).
			LEFT_JOIN(Requests, Requests.WorkflowVersionID.EQ(WorkflowVersions.ID)),
	).WHERE(Workflows.ID.EQ(postgres.Int32(workflowId)))

	workflow := results.Workflow{}

	err := statement.QueryContext(ctx, db, &workflow)

	return workflow, err
}

func (r *WorkflowRepository) UpdateWorkflow(ctx context.Context, tx *sql.Tx, workflow model.Workflows) error {
	Workflows := table.Workflows

	workflow.UpdatedAt = time.Now().UTC().Add(7 * time.Hour)

	columns := Workflows.AllColumns.Except(Workflows.ID, Workflows.CreatedAt, Workflows.DeletedAt)

	statement := Workflows.UPDATE(columns).MODEL(workflow).WHERE(Workflows.ID.EQ(postgres.Int32(workflow.ID))).RETURNING(Workflows.ID)

	err := statement.QueryContext(ctx, tx, &workflow)

	return err
}
