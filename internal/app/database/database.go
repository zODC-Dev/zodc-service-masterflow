package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/configs"
)

func ConnectDatabase() *sql.DB {
	db, err := sql.Open("pgx", configs.Env.DB_DSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return db
}
