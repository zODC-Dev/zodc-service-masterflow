package main

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/configs"
	db "github.com/zODC-Dev/zodc-service-masterflow/internal/app/database"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/externals"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/routes"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/services"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/nats"
)

// @title zODC Masterflow Service API
// @version 1.0
// @description This is the API documentation for the zODC Masterflow Service. It provides endpoints for managing workflows, requests, forms, and nodes.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1  // Adjusted BasePath for clarity (e.g., /api/v1) - Ensure routes match this prefix

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	ctx := context.Background()
	e := echo.New()
	{
		// e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"*"},
			AllowHeaders: []string{"*"},
		}))
	}

	//Database Setup
	database := db.ConnectDatabase()

	//Nats Setup
	natsSubscriber := setupNats(ctx, database)
	// Đảm bảo cleanup khi chương trình kết thúc
	if natsSubscriber != nil {
		defer func() {
			ctx := context.Background()
			if err := natsSubscriber.Shutdown(ctx); err != nil {
				slog.Error("Error shutting down NATS subscriber service", "error", err)
			}
		}()
	}

	//Route Setup
	routes.RegisterRoutes(e, database)

	slog.Error(e.Start(configs.Env.SERVER_ADDRESS).Error())
}

func setupNats(ctx context.Context, database *sql.DB) *services.NatsSubscriberService {
	// Initialize repositories
	nodeRepo := repositories.NewNodeRepository()
	requestRepo := repositories.NewRequestRepository()
	connectionRepo := repositories.NewConnectionRepository()
	formRepo := repositories.NewFormRepository()
	historyRepo := repositories.NewHistoryRepository()

	// Initialize NATS client
	natsConfig := nats.DefaultConfig()
	natsConfig.URL = configs.Env.NATS_URL // Ensure NATS_URL is defined in your .env file
	natsClient, err := nats.NewNATSClient(natsConfig)
	if err != nil {
		slog.Error("Failed to create NATS client", "error", err)
		return nil
	}

	// Khởi tạo các services cần thiết
	// Set up UserAPI if needed for NodeService
	userApi := externals.NewUserAPI()

	notificationService := services.NewNotificationService(database, natsClient, userApi, requestRepo)
	historyService := services.NewHistoryService(database, historyRepo, userApi)
	formService := services.NewFormService(database, formRepo, natsClient)

	// Create NatsService needed for NodeService
	natsService := services.NewNatsService(services.NatsService{
		NatsClient:  natsClient,
		NodeRepo:    nodeRepo,
		RequestRepo: requestRepo,
		FormRepo:    formRepo,
	})

	// Khởi tạo RequestService - cần cho NodeService
	requestService := services.NewRequestService(services.RequestService{
		DB:             database,
		RequestRepo:    requestRepo,
		NodeRepo:       nodeRepo,
		ConnectionRepo: connectionRepo,
		HistoryService: historyService,
	})

	// Initialize the NodeService with all dependencies
	nodeService := services.NewNodeService(services.NodeService{
		DB:                  database,
		NodeRepo:            nodeRepo,
		ConnectionRepo:      connectionRepo,
		RequestRepo:         requestRepo,
		FormRepo:            formRepo,
		NatsClient:          natsClient,
		NatsService:         natsService,
		UserAPI:             userApi,
		HistoryService:      historyService,
		NotificationService: notificationService,
		RequestService:      requestService,
		FormService:         formService,
	})

	// Initialize and start NATS subscriber service
	natsSubscriberService := services.NewNatsSubscriberService(
		natsClient,
		database,
		nodeRepo,
		requestRepo,
		nodeService,
	)

	// Start the subscriber service in a goroutine
	go func() {
		if err := natsSubscriberService.Start(ctx); err != nil {
			slog.Error("Failed to start NATS subscriber service", "error", err)
		}
	}()

	slog.Info("NATS subscriber service started")
	return natsSubscriberService
}
