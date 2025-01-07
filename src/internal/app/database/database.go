package db

import (
	"log/slog"

	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/configs"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	dsn := configs.Env.DATABASE_POSTGRE_DSN
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	defer slog.Info("PostgreSQL running")

	if err != nil {
		return nil, err
	}

	//Auto Migration
	db.AutoMigrate(&entities.FormExcel{}, &entities.Form{})

	return db, nil
}
