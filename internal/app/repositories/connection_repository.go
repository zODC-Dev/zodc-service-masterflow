package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/table"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/results"
)

type ConnectionRepository struct{}

func NewConnectionRepository() *ConnectionRepository {
	return &ConnectionRepository{}
}

func (r *ConnectionRepository) FindConnectionsByToNodeId(ctx context.Context, db *sql.DB, toNodeId string) ([]model.Connections, error) {
	Connections := table.Connections

	statement := postgres.SELECT(Connections.AllColumns).FROM(Connections).WHERE(Connections.ToNodeID.EQ(postgres.String(toNodeId)))

	result := []model.Connections{}
	err := statement.QueryContext(ctx, db, &result)

	return result, err
}

func (r *ConnectionRepository) FindConnectionsWithToNodesByFromNodeId(ctx context.Context, db *sql.DB, fromNodeId string) ([]results.ConnectionWithNode, error) {
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

func (r *ConnectionRepository) UpdateConnection(ctx context.Context, tx *sql.Tx, connection model.Connections) error {
	Connections := table.Connections

	connection.UpdatedAt = time.Now()

	columns := Connections.AllColumns.Except(Connections.ID, Connections.CreatedAt, Connections.DeletedAt)

	statment := Connections.UPDATE(columns).MODEL(connection)

	err := statment.QueryContext(ctx, tx, &connection)

	return err
}

func (r *ConnectionRepository) FindAllConnectionByRequestId(ctx context.Context, db *sql.DB, requestId int32) ([]model.Connections, error) {
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

func (r *ConnectionRepository) CreateConnections(ctx context.Context, tx *sql.Tx, connections []model.Connections) error {
	Connections := table.Connections

	columns := Connections.AllColumns.Except(Connections.CreatedAt, Connections.UpdatedAt, Connections.DeletedAt)

	statement := Connections.INSERT(columns).MODELS(connections)

	err := statement.QueryContext(ctx, tx, &connections)

	if err != nil {
		fmt.Println(statement.DebugSql())
	}

	return err
}
