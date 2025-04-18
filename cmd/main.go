package main

import (
	"context"
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
	// Create context for services
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

	// Initialize repositories
	nodeRepo := repositories.NewNodeRepository()
	requestRepo := repositories.NewRequestRepository()

	// Initialize NATS client
	natsConfig := nats.DefaultConfig()
	natsConfig.URL = configs.Env.NATS_URL // Ensure NATS_URL is defined in your .env file
	natsClient, err := nats.NewNATSClient(natsConfig)
	if err != nil {
		slog.Error("Failed to create NATS client", "error", err)
	} else {
		// Initialize and start NATS subscriber service
		natsSubscriberService := services.NewNatsSubscriberService(
			natsClient,
			database,
			nodeRepo,
			requestRepo,
			nil, // We will create a proper NodeService below and set it
		)

		// Initialize other required dependencies for NodeService
		connectionRepo := repositories.NewConnectionRepository()
		formRepo := repositories.NewFormRepository()

		// Set up UserAPI if needed for NodeService
		userApi := externals.NewUserAPI()

		// Create NatsService needed for NodeService
		natsService := services.NewNatsService(services.NatsService{
			NatsClient:  natsClient,
			NodeRepo:    nodeRepo,
			RequestRepo: requestRepo,
			FormRepo:    formRepo,
		})

		// Initialize the NodeService with all dependencies
		nodeService := services.NewNodeService(services.NodeService{
			DB:             database,
			NodeRepo:       nodeRepo,
			ConnectionRepo: connectionRepo,
			RequestRepo:    requestRepo,
			FormRepo:       formRepo,
			NatsClient:     natsClient,
			NatsService:    natsService,
			UserAPI:        userApi,
		})

		// Now set the NodeService in the NatsSubscriberService
		natsSubscriberService.NodeService = nodeService

		// Start the subscriber service in a goroutine
		go func() {
			if err := natsSubscriberService.Start(ctx); err != nil {
				slog.Error("Failed to start NATS subscriber service", "error", err)
			}
		}()

		// Ensure graceful shutdown
		defer func() {
			if err := natsSubscriberService.Shutdown(ctx); err != nil {
				slog.Error("Error shutting down NATS subscriber service", "error", err)
			}
		}()

		slog.Info("NATS subscriber service started")
	}

	//Route Setup
	routes.RegisterRoutes(e, database)

	slog.Error(e.Start(configs.Env.SERVER_ADDRESS).Error())
}
