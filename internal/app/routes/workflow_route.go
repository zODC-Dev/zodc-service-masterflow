package routes

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/controllers"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/services"
)

func WorkflowRoute(group *echo.Group, db *sql.DB) {
	nodeRepo := repositories.NewNodeRepository()
	workflowRepo := repositories.NewWorkflowRepository()
	nodeConnectionRepo := repositories.NewNodeConnectionRepository()

	workflowService := services.NewWorkflowService(db, nodeRepo, workflowRepo, nodeConnectionRepo)

	workflowController := controllers.NewWorkflowController(workflowService)

	workflowRoute := group.Group("/workflows")
	{
		workflowRoute.POST("/create", workflowController.Create)
		workflowRoute.GET("/all", workflowController.FindAll)
	}
}
