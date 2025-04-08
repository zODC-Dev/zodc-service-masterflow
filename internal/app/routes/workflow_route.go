package routes

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/controllers"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/externals"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/nats"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/services"
)

func WorkflowRoute(group *echo.Group, db *sql.DB) {
	// Repositories
	workflowRepo := repositories.NewWorkflowRepository()
	formRepo := repositories.NewFormRepository()
	categoryRepo := repositories.NewCategoryRepository()
	requestRepo := repositories.NewRequestRepository()
	connectionRepo := repositories.NewConnectionRepository()
	nodeRepo := repositories.NewNodeRepository()
	natsClient := nats.GetNATSClient()

	// Apis
	userApi := externals.NewUserAPI()

	// External Services
	nodeService := services.NewNodeService(services.NodeService{
		DB:             db,
		NodeRepo:       nodeRepo,
		ConnectionRepo: connectionRepo,
		RequestRepo:    requestRepo,
	})

	natsService := services.NewNatsService(services.NatsService{
		NodeRepo:   nodeRepo,
		NatsClient: natsClient,
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
		NatsClient:     natsClient,
		NatsService:    natsService,
	})

	nodeService.WorkflowService = workflowService

	workflowController := controllers.NewWorkflowController(workflowService)

	workflowRoute := group.Group("/workflows")
	{
		workflowRoute.GET("", workflowController.FindAllWorkflow)
		workflowRoute.POST("/create", workflowController.CreateWorkflow)
		workflowRoute.GET("/:id", workflowController.FindOneWorkflowDetail)

		workflowRoute.POST("/start", workflowController.StartWorkflow)
		workflowRoute.PUT("/:id/archive", workflowController.ArchiveWorkflow)

	}
}
