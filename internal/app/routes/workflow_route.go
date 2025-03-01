package routes

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/controllers"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/services"
)

func WorkflowRoute(group *echo.Group, db *sql.DB) {
	workflowRepo := repositories.NewWorkflowRepository()
	formRepo := repositories.NewFormRepository()
	categoryRepo := repositories.NewCategoryRepository()

	workflowService := services.NewWorkflowService(db, workflowRepo, formRepo, categoryRepo)
	workflowController := controllers.NewWorkflowController(workflowService)

	workflowRoute := group.Group("/workflows")
	{
		workflowRoute.GET("", workflowController.FindAllWorkflow)
		workflowRoute.POST("/create", workflowController.CreateWorkflow)
		workflowRoute.GET("/:id", workflowController.FindOneWorkflowDetail)
	}
}
