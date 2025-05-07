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

func RequestRoute(group *echo.Group, db *sql.DB) {
	requestRepo := repositories.NewRequestRepository()
	workflowRepo := repositories.NewWorkflowRepository()
	formRepo := repositories.NewFormRepository()
	categoryRepo := repositories.NewCategoryRepository()
	connectionRepo := repositories.NewConnectionRepository()
	nodeRepo := repositories.NewNodeRepository()
	historyRepo := repositories.NewHistoryRepository()

	natsClient := nats.GetNATSClient()
	userApi := externals.NewUserAPI()

	formService := services.NewFormService(db, formRepo, natsClient)
	historyService := services.NewHistoryService(db, historyRepo, userApi)
	notificationService := services.NewNotificationService(db, natsClient, userApi, requestRepo, historyService)

	nodeService := services.NewNodeService(services.NodeService{
		DB:                  db,
		NodeRepo:            nodeRepo,
		FormRepo:            formRepo,
		HistoryService:      historyService,
		NotificationService: notificationService,
	})

	natsService := services.NewNatsService(services.NatsService{
		NatsClient:  natsClient,
		NodeRepo:    nodeRepo,
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
		HistoryService:      historyService,
		NotificationService: notificationService,
	})

	requestService := services.NewRequestService(services.RequestService{
		DB:                  db,
		RequestRepo:         requestRepo,
		UserAPI:             userApi,
		WorkflowRepo:        workflowRepo,
		WorkflowService:     workflowService,
		ConnectionRepo:      connectionRepo,
		NodeRepo:            nodeRepo,
		NatsService:         natsService,
		NodeService:         nodeService,
		FormService:         formService,
		FormRepo:            formRepo,
		HistoryRepo:         historyRepo,
		HistoryService:      historyService,
		NotificationService: notificationService,
	})

	nodeService.RequestService = requestService

	requestController := controllers.NewRequestController(requestService)

	requestRoute := group.Group("/requests")
	{
		requestRoute.GET("", requestController.FindAllRequest)
		requestRoute.GET("/count", requestController.GetRequestCount)
		requestRoute.GET("/detail/:id", requestController.GetRequestDetail)
		requestRoute.GET("/tasks/:id", requestController.GetRequestTasks)

		requestRoute.GET("/tasks/projects", requestController.GetRequestTasksByProject)

		requestRoute.GET("/overview/:id", requestController.GetRequestOverview)

		requestRoute.GET("/:id/sub-requests", requestController.FindAllSubRequestByRequestId)

		requestRoute.PUT("/:id", requestController.UpdateRequest)

		requestRoute.GET("/tasks/count", requestController.GetRequestTasksCount)

		requestRoute.GET("/:id/completed-form", requestController.GetRequestCompletedForm)
		requestRoute.GET("/:id/completed-form/:dataId", requestController.GetRequestCompletedFormApproval)

		requestRoute.GET("/:id/file-manager", requestController.GetRequestFileManager)

		requestRoute.GET("/:id/history", requestController.FindAllHistoryByRequestId)

		requestRoute.GET("/report/mid-sprint-tasks", requestController.ReportMidSprintTasks)

		requestRoute.PUT("/:id/cancel", requestController.CancelRequest)

		requestRoute.GET("/retrospective/report/:sprintId", requestController.GetRetrospectiveReport)

		requestRoute.POST("/:id/complete", requestController.CompleteRequestHandler)

		requestRoute.POST("/sprint/:sprintId/complete", requestController.CompleteAllRequestBySprintIdHandler)
	}
}
