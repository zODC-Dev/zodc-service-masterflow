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

func NodeRoute(group *echo.Group, db *sql.DB) {
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

	natsService := services.NewNatsService(services.NatsService{
		NatsClient:  natsClient,
		NodeRepo:    nodeRepo,
		RequestRepo: requestRepo,
	})

	nodeService := services.NewNodeService(services.NodeService{
		NodeRepo:       nodeRepo,
		ConnectionRepo: connectionRepo,
		RequestRepo:    requestRepo,
		DB:             db,
		FormRepo:       formRepo,
		NatsClient:     natsClient,
		NatsService:    natsService,
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
		NatsService:    natsService,
		NatsClient:     natsClient,
	})

	nodeService.WorkflowService = workflowService

	nodeController := controllers.NewNodeController(nodeService)

	nodeRoute := group.Group("/nodes")
	{
		nodeRoute.POST("/:id/complete", nodeController.CompleteNode)
		nodeRoute.POST("/:id/start", nodeController.StartNode)
		nodeRoute.POST("/:id/approve", nodeController.ApproveNode)
		nodeRoute.POST("/:id/reject", nodeController.RejectNode)

		nodeRoute.POST("/:id/reassign/:userId", nodeController.ReassignNode)

		nodeRoute.POST("/:id/forms/:formId/submit", nodeController.SubmitNodeForm)
		nodeRoute.POST("/:id/forms/:formId/approve", nodeController.ApproveNodeForm)
		nodeRoute.POST("/:id/forms/:formId/reject", nodeController.RejectNodeForm)

		nodeRoute.GET("/:id/forms/:permission", nodeController.GetNodeFormWithPermission)
		nodeRoute.GET("/:id/jira-form", nodeController.GetNodeJiraForm)

		nodeRoute.GET("/:id/task", nodeController.GetNodeTaskDetail)

		nodeRoute.GET("/stories", nodeController.GetNodeStoryByAssignee)
	}
}
