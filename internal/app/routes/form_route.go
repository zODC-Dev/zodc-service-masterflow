package routes

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/controllers"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/nats"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/services"
)

func FormRoute(group *echo.Group, db *sql.DB) {
	formRepo := repositories.NewFormRepository()
	natsClient := nats.GetNATSClient()
	formService := services.NewFormService(db, formRepo, natsClient)

	formController := controllers.NewFormController(formService)

	formRoute := group.Group("/forms")
	{
		formRoute.POST("/templates/create", formController.CreateFormTemplate)
		formRoute.GET("/templates/all", formController.FindAllFormTemplate)

		formRoute.PUT("/templates/:formTemplateId/edit", formController.UpdateFormTemplate)
		formRoute.PUT("/templates/:formTemplateId/config", formController.ConfigFormTemplate)

		formRoute.GET("/templates/:formTemplateId/detail", formController.FindOneFormTemplateDetailByFormTemplateId)

		formRoute.GET("/data/:dataId", formController.FindOneFormDataByFormDataId)

		formRoute.GET("/templates/edit-profile", formController.GetEditProfileFormTemplate)
		formRoute.GET("/templates/performance-evaluate", formController.GetPerformanceEvaluateFormTemplate)
	}
}
