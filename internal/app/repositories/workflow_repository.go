package repositories

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/table"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/filters"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
)

type WorkflowRepository struct{}

func NewWorkflowRepository() *WorkflowRepository {
	return &WorkflowRepository{}
}

func (r *WorkflowRepository) Create(ctx context.Context, tx *sql.Tx, workflow model.Workflows) (model.Workflows, error) {
	Workflows := table.Workflows

	workflowInsertColumns := Workflows.AllColumns.Except(Workflows.ID, Workflows.CreatedAt, Workflows.UpdatedAt, Workflows.DeletedAt)

	workflowStmt := Workflows.INSERT(workflowInsertColumns).MODEL(workflow).RETURNING(Workflows.ID)

	if err := workflowStmt.QueryContext(ctx, tx, &workflow); err != nil {
		return workflow, err
	}

	return workflow, nil
}

func (r *WorkflowRepository) FindAll(ctx context.Context, db *sql.DB, workflowFilter filters.WorkflowFilter) ([]types.WorkflowType, error) {
	Workflows := table.Workflows
	Nodes := table.Nodes
	NodeConnections := table.NodeConnections
	NodeGroups := table.NodeGroups
	Categories := table.Categories

	stmt := postgres.SELECT(
		Workflows.AllColumns,
		Nodes.AllColumns,
		NodeConnections.AllColumns,
		NodeGroups.AllColumns,
		Categories.AllColumns,
	).FROM(
		Workflows.
			LEFT_JOIN(Nodes, Workflows.ID.EQ(Nodes.WorkflowID)).
			LEFT_JOIN(NodeConnections, Workflows.ID.EQ(NodeConnections.WorkflowID)).
			LEFT_JOIN(NodeGroups, Workflows.ID.EQ(NodeGroups.WorkflowID)).
			LEFT_JOIN(Categories, Workflows.CategoryID.EQ(Categories.ID)),
	)

	if workflowFilter.CategoryID != "" {
		categoryIdInt, err := strconv.Atoi(workflowFilter.CategoryID)
		if err != nil {
			return []types.WorkflowType{}, err
		}

		stmt.WHERE(Workflows.CategoryID.EQ(postgres.Int(int64(categoryIdInt))))
	}

	if workflowFilter.Type != "" {
		stmt.WHERE(Workflows.Type.EQ(postgres.String(workflowFilter.Type)))
	}

	if workflowFilter.Search != "" {
		stmt.WHERE(postgres.LOWER(Workflows.Title).LIKE(postgres.LOWER(postgres.String("%" + workflowFilter.Search + "%"))))
	}

	workflows := []types.WorkflowType{}
	err := stmt.QueryContext(ctx, db, &workflows)

	for i := range workflows {
		if workflows[i].Nodes == nil {
			workflows[i].Nodes = []model.Nodes{}
		}
		if workflows[i].Groups == nil {
			workflows[i].Groups = []model.NodeGroups{}
		}
		if workflows[i].Connections == nil {
			workflows[i].Connections = []model.NodeConnections{}
		}
	}

	return workflows, err
}
