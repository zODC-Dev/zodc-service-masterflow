package db

import (
	"log/slog"
	"os"

	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/configs"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dsn := configs.Env.DATABASE_POSTGRE_DSN
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	defer slog.Info("Database running")

	if err != nil {
		slog.Error("Database fail: ", slog.Any("error", err))

		//End App
		os.Exit(1)
	}

	//Auto Migration
	db.AutoMigrate(&entities.FormExcel{}, &entities.Form{})

	return db
}
