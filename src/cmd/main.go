package main

import (
	"net/http"

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
		app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			Skipper:      middleware.DefaultSkipper,
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		}))
	}

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
