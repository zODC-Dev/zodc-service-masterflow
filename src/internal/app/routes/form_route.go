package routes

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/controllers"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/services"
)

func FormRoute(group *echo.Group, db *sql.DB) {
	formRepo := repositories.NewFormRepository(db)
	formService := services.NewFormService(formRepo)
	formController := controllers.NewFormController(formService)

	formRoute := group.Group("/forms")
	{
		formRoute.POST("/create", formController.Create)
		formRoute.GET("/all", formController.FindAll)
	}
}
