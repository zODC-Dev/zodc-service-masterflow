package main

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/configs"
	db "github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/database"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/routes"
)

func main() {
	app := echo.New()
	{
		app.Use(middleware.Logger())
		app.Use(middleware.Recover())
	}

	//Database Setup
	db, err := db.ConnectDatabase()
	if err != nil {
		slog.Error("Database fail: ", slog.Any("error", err))
		return
	}

	//Route Setup
	routeGroup := app.Group("/api/v1")
	{
		routes.FormRoute(routeGroup, db)
		routes.UtilRoute(routeGroup, db)
	}

	app.Logger.Fatal(app.Start(configs.Env.SERVER_ADDRESS))
}
