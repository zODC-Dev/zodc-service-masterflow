package routes

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/controllers"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/externals"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/services"
)

func NodeRoute(group *echo.Group, db *sql.DB) {
	// Repositories
	workflowRepo := repositories.NewWorkflowRepository()
	formRepo := repositories.NewFormRepository()
	categoryRepo := repositories.NewCategoryRepository()
	requestRepo := repositories.NewRequestRepository()
	connectionRepo := repositories.NewConnectionRepository()
	nodeRepo := repositories.NewNodeRepository()

	// Apis
	userApi := externals.NewUserAPI()

	//

	nodeService := services.NewNodeService(services.NodeService{
		NodeRepo:       nodeRepo,
		ConnectionRepo: connectionRepo,
		RequestRepo:    requestRepo,
	})

	workflowService := services.NewWorkflowService(services.WorkflowService{
		DB:             db,
		WorkflowRepo:   workflowRepo,
		FormRepo:       formRepo,
		CategoryRepo:   categoryRepo,
		UserAPI:        userApi,
		RequestRepo:    requestRepo,
		ConnectionRepo: connectionRepo,
		NodeRepo:       nodeRepo,
		NodeService:    nodeService,
	})

	nodeService.WorkflowService = workflowService

	nodeController := controllers.NewNodeController(nodeService)

	nodeRoute := group.Group("/nodes")
	{
		nodeRoute.POST("/:id/complete", nodeController.CompleteNode)
		nodeRoute.POST("/:id/start", nodeController.StartNode)
	}
}
