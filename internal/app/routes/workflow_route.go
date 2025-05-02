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
	historyRepo := repositories.NewHistoryRepository()

	// Nats
	natsClient := nats.GetNATSClient()

	// Apis
	userApi := externals.NewUserAPI()

	// External Services
	historyService := services.NewHistoryService(db, historyRepo, userApi)
	notificationService := services.NewNotificationService(db, natsClient, userApi, requestRepo, historyService)

	requestService := services.NewRequestService(services.RequestService{
		DB:             db,
		RequestRepo:    requestRepo,
		UserAPI:        userApi,
		ConnectionRepo: connectionRepo,
		NodeRepo:       nodeRepo,
		FormRepo:       formRepo,
		HistoryRepo:    historyRepo,
		HistoryService: historyService,
	})

	nodeService := services.NewNodeService(services.NodeService{
		DB:             db,
		NodeRepo:       nodeRepo,
		ConnectionRepo: connectionRepo,
		RequestRepo:    requestRepo,
		FormRepo:       formRepo,
		HistoryService: historyService,
		RequestService: requestService,
	})

	natsService := services.NewNatsService(services.NatsService{
		NodeRepo:    nodeRepo,
		NatsClient:  natsClient,
		RequestRepo: requestRepo,
		FormRepo:    formRepo,
	})

	workflowService := services.NewWorkflowService(services.WorkflowService{
		DB:                  db,
		WorkflowRepo:        workflowRepo,
		FormRepo:            formRepo,
		CategoryRepo:        categoryRepo,
		UserAPI:             userApi,
		RequestRepo:         requestRepo,
		ConnectionRepo:      connectionRepo,
		NodeRepo:            nodeRepo,
		NodeService:         nodeService,
		NatsClient:          natsClient,
		NatsService:         natsService,
		NotificationService: notificationService,
		HistoryService:      historyService,
	})

	nodeService.WorkflowService = workflowService

	workflowController := controllers.NewWorkflowController(workflowService)

	workflowRoute := group.Group("/workflows")
	{
		workflowRoute.GET("", workflowController.FindAllWorkflow)
		workflowRoute.POST("/create", workflowController.CreateWorkflow)

		workflowRoute.GET("/:id", workflowController.FindOneWorkflowDetail)
		workflowRoute.PUT("/:id/edit", workflowController.UpdateWorkflow)

		workflowRoute.POST("/start", workflowController.StartWorkflow)
		workflowRoute.PUT("/:id/archive", workflowController.ArchiveWorkflow)

	}
}
