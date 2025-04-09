package routes

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/controllers"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/externals"
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
	nodeService := services.NewNodeService(services.NodeService{
		DB:       db,
		NodeRepo: nodeRepo,
		FormRepo: formRepo,
	})

	userApi := externals.NewUserAPI()

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

	requestService := services.NewRequestService(services.RequestService{
		DB:              db,
		RequestRepo:     requestRepo,
		UserAPI:         userApi,
		WorkflowService: workflowService,
	})

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

	}

}
