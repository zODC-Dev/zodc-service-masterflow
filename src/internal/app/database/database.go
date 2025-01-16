package db

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

func ConnectDatabase() *pgx.Conn {
	ctx := context.Background()
	db, err := pgx.Connect(ctx, "postgres://zodcdbuser:D4yl5m4tkh4usi3um4nh@thisisdbformasterzodcserviceflowonly.thanhf.dev:5432/zodc_masterflow")
	if err != nil {
		slog.Error("Database fail: ", slog.Any("error", err))
		os.Exit(1)
	}

	return db
}
