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

type HistoryRepository struct{}

func NewHistoryRepository() *HistoryRepository {
	return &HistoryRepository{}
}

func (r *HistoryRepository) CreateHistory(ctx context.Context, tx *sql.Tx, history model.Histories) error {
	Histories := table.Histories

	columns := Histories.AllColumns.Except(Histories.ID, Histories.CreatedAt, Histories.UpdatedAt, Histories.DeletedAt)

	statement := Histories.INSERT(columns).MODEL(history).RETURNING(Histories.ID)

	err := statement.QueryContext(ctx, tx, &history)

	return err
}

func (r *HistoryRepository) FindAllHistoryByRequestId(ctx context.Context, db *sql.DB, requestId int32) ([]results.HistoryResult, error) {
	Histories := table.Histories
	Nodes := table.Nodes

	statement := postgres.SELECT(
		Histories.AllColumns,
		Nodes.AllColumns,
	).FROM(
		Histories.
			LEFT_JOIN(Nodes, Nodes.ID.EQ(Histories.NodeID)),
	).WHERE(
		Histories.RequestID.EQ(postgres.Int32(requestId)),
	).ORDER_BY(
		Histories.CreatedAt.DESC(),
		Histories.TypeAction.EQ(postgres.String(constants.HistoryTypeNewTask)).DESC(),
	)

	var histories []results.HistoryResult
	err := statement.QueryContext(ctx, db, &histories)

	return histories, err
}
