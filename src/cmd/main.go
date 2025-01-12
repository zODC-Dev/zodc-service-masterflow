package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/configs"
	db "github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/database"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/routes"
)

func main() {
	app := echo.New()
	{
		app.Use(middleware.Logger())
		app.Use(middleware.Recover())
		app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"*"},
			AllowHeaders: []string{"*"},
		}))
	}

	//Setup swagger
	app.GET("/docs/*", echoSwagger.WrapHandler)

	//Database Setup
	db := db.ConnectDatabase()

	//Route Setup
	routeGroup := app.Group(configs.Server.API_Prefix)
	{
		routes.FormRoute(routeGroup, db)
		routes.UtilRoute(routeGroup, db)
	}

	app.Logger.Fatal(app.Start(configs.Env.SERVER_ADDRESS))
}
