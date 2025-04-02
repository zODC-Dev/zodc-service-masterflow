package routes

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/controllers"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/externals"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/services"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/nats"
)

func WorkflowRoute(group *echo.Group, db *sql.DB) {
	// Repositories
	workflowRepo := repositories.NewWorkflowRepository()
	formRepo := repositories.NewFormRepository()
	categoryRepo := repositories.NewCategoryRepository()
	requestRepo := repositories.NewRequestRepository()
	connectionRepo := repositories.NewConnectionRepository()
	nodeRepo := repositories.NewNodeRepository()
	natsClient, err := nats.NewNATSClient(nats.DefaultConfig())
	if err != nil {
		slog.Error("Failed to create NATS client", "error", err)
		os.Exit(1)
	}

	// Apis
	userApi := externals.NewUserAPI()

	// External Services
	nodeService := services.NewNodeService(services.NodeService{
		DB:             db,
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
		NatsClient:     natsClient,
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
