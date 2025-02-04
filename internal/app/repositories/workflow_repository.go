package repositories

import (
	"context"
	"database/sql"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/model"
	. "github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/table"
)

type WorkflowRepository struct{}

func NewWorkflowRepository() *WorkflowRepository {
	return &WorkflowRepository{}
}

func (r *WorkflowRepository) Create(ctx context.Context, tx *sql.Tx, workflow model.Workflows) (model.Workflows, error) {
	workflowInsertColumns := Workflows.AllColumns.Except(Workflows.ID, Workflows.CreatedAt, Workflows.UpdatedAt, Workflows.DeletedAt)

	workflowStmt := Workflows.INSERT(workflowInsertColumns).MODEL(workflow).RETURNING(Workflows.ID)

	if err := workflowStmt.QueryContext(ctx, tx, &workflow); err != nil {
		return workflow, err
	}

	return workflow, nil
}
