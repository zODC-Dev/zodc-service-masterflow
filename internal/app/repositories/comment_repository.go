package repositories

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/table"
)

type CommentRepository struct{}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{}
}

func (r *CommentRepository) CreateComment(ctx context.Context, tx *sql.Tx, comment model.Comments) error {
	Comments := table.Comments

	columns := Comments.AllColumns.Except(Comments.ID, Comments.CreatedAt, Comments.UpdatedAt, Comments.DeletedAt)

	statement := Comments.INSERT(columns).MODEL(comment)

	return statement.QueryContext(ctx, tx, &comment)
}

func (r *CommentRepository) FindAllCommentByNodeId(ctx context.Context, db *sql.DB, nodeId string) ([]model.Comments, error) {
	Comments := table.Comments
	Nodes := table.Nodes

	statement := postgres.SELECT(
		Comments.AllColumns,
	).FROM(
		Comments.
			LEFT_JOIN(Nodes, Nodes.ID.EQ(Comments.NodeID)),
	).WHERE(
		Nodes.ID.EQ(postgres.String(nodeId)),
	)

	result := []model.Comments{}
	err := statement.QueryContext(ctx, db, &result)

	return result, err

}
