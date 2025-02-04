package routes

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/controllers"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/services"
)

func FormRoute(group *echo.Group, db *sql.DB) {
	formRepo := repositories.NewFormRepository()
	formFieldRepo := repositories.NewFormFieldRepository()

	formService := services.NewFormService(db, formRepo, formFieldRepo)

	formController := controllers.NewFormController(formService)

	formRoute := group.Group("/forms")
	{
		formRoute.POST("/create", formController.Create)
		formRoute.GET("/all", formController.FindAll)
	}
}
