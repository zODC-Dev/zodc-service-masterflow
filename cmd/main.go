package main

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/configs"
	db "github.com/zODC-Dev/zodc-service-masterflow/internal/app/database"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/routes"
)

// @title zODC Masterflow Service API
// @version 1.0
// @description This is the API documentation for the zODC Masterflow Service. It provides endpoints for managing workflows, requests, forms, and nodes.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1  // Adjusted BasePath for clarity (e.g., /api/v1) - Ensure routes match this prefix

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
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
