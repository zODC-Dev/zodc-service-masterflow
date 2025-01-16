package routes

import (
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/controllers"
	database "github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/database/generated"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/services"
)

func FormRoute(group *echo.Group, db *pgx.Conn, queries *database.Queries) {
	formService := services.NewFormService(db, queries)
	formController := controllers.NewFormController(formService)

	formRoute := group.Group("/forms")
	{
		formRoute.POST("/create", formController.Create)
		formRoute.GET("/all", formController.FindAll)
	}
}
