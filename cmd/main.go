package main

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/configs"
	db "github.com/zODC-Dev/zodc-service-masterflow/internal/app/database"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/routes"
)

// @title Echo Swagger API
// @version 1.0
// @host localhost:8080 // PHẢI KHỚP PORT
// @BasePath /
func main() {
	e := echo.New()
	{
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"*"},
			AllowHeaders: []string{"*"},
		}))
	}

	//Database Setup
	db := db.ConnectDatabase()

	//Route Setup
	routes.RegisterRoutes(e, db)

	slog.Error(e.Start(configs.Env.SERVER_ADDRESS).Error())
}
