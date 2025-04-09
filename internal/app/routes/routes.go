package routes

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/middlewares"
)

func RegisterRoutes(e *echo.Echo, db *sql.DB) {

	//Swagger Setup
	RegisterSwaggerRoute(e)

	// Api V1 Group
	apiV1Group := e.Group("/api/v1")

	// Middleware Setup
	apiV1Group.Use(middlewares.ExtractUserMiddleware())
	{
		FormRoute(apiV1Group, db)
		WorkflowRoute(apiV1Group, db)
		UtilRoute(apiV1Group)
		CategoryRoute(apiV1Group, db)
		RequestRoute(apiV1Group, db)
		NodeRoute(apiV1Group, db)
	}
}
