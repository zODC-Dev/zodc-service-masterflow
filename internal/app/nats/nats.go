package nats

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"

	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/configs"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/externals"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/services"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/nats"
)

var (
	natsClient     *nats.NATSClient
	natsClientOnce sync.Once
)

// GetNATSClient trả về instance của NATS client (singleton pattern)
func GetNATSClient() *nats.NATSClient {
	natsClientOnce.Do(func() {
		natsConfig := nats.DefaultConfig()
		natsConfig.URL = configs.Env.NATS_URL

		var err error
		natsClient, err = nats.NewNATSClient(natsConfig)
		if err != nil {
			slog.Error("Failed to connect to NATS", "error", err)
			natsClient = nil
		} else {
			slog.Info("Connected to NATS server", "url", natsConfig.URL)
		}
	})

	return natsClient
}

func SetupNats(database *sql.DB) *services.NatsSubscriberService {
	ctx := context.Background()

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

	historyService := services.NewHistoryService(database, historyRepo, userApi)
	notificationService := services.NewNotificationService(database, natsClient, userApi, requestRepo, historyService)

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

	if natsSubscriberService != nil {
		if err := natsSubscriberService.Shutdown(ctx); err != nil {
			slog.Error("Error shutting down NATS subscriber service", "error", err)
		}
	}

	slog.Info("NATS subscriber service started")
	return natsSubscriberService
}
