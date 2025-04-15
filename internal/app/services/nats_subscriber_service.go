package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	natslib "github.com/nats-io/nats.go"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/nats"
)

// NatsSubscriberService handles listening for messages from NATS
type NatsSubscriberService struct {
	NatsClient  *nats.NATSClient
	DB          *sql.DB
	NodeRepo    *repositories.NodeRepository
	RequestRepo *repositories.RequestRepository
	// Store subscriptions to properly unsubscribe later
	subscriptions []*natslib.Subscription
}

// NewNatsSubscriberService creates a new instance of NatsSubscriberService
func NewNatsSubscriberService(natsClient *nats.NATSClient, db *sql.DB, nodeRepo *repositories.NodeRepository, requestRepo *repositories.RequestRepository) *NatsSubscriberService {
	return &NatsSubscriberService{
		NatsClient:    natsClient,
		DB:            db,
		NodeRepo:      nodeRepo,
		RequestRepo:   requestRepo,
		subscriptions: make([]*natslib.Subscription, 0),
	}
}

// Start initializes and starts the subscriber service
func (s *NatsSubscriberService) Start(ctx context.Context) error {
	slog.Info("Starting NATS Subscriber Service")

	// Subscribe to topics
	if err := s.subscribeToTopics(ctx); err != nil {
		return fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	slog.Info("NATS Subscriber Service started successfully")
	return nil
}

// subscribeToTopics registers subscriptions for each topic
func (s *NatsSubscriberService) subscribeToTopics(ctx context.Context) error {
	// List of topics to subscribe to
	topics := []string{
		// Add your topics here
		"topic1.event",
		"topic2.event",
	}

	for _, topic := range topics {
		slog.Info("Subscribing to topic", "topic", topic)

		// Use the Subscribe method from pkg/nats/subscribe.go
		subscription, err := s.NatsClient.Subscribe(topic, s.handleMessage)
		if err != nil {
			return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
		}

		// Store subscription for later cleanup
		s.subscriptions = append(s.subscriptions, subscription)
	}

	return nil
}

// handleMessage processes received messages
func (s *NatsSubscriberService) handleMessage(msg *natslib.Msg) {
	slog.Info("Received message", "subject", msg.Subject, "data_length", len(msg.Data))

	// Start transaction
	tx, err := s.DB.Begin()
	if err != nil {
		slog.Error("Failed to start transaction", "error", err)
		return
	}
	defer tx.Rollback()

	// Process message based on subject
	var processErr error
	switch msg.Subject {
	case "topic1.event":
		processErr = s.handleTopic1(tx, msg.Data)
	case "topic2.event":
		processErr = s.handleTopic2(tx, msg.Data)
	default:
		slog.Warn("No handler for subject", "subject", msg.Subject)
	}

	if processErr != nil {
		slog.Error("Failed to process message", "subject", msg.Subject, "error", processErr)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		slog.Error("Failed to commit transaction", "error", err)
		return
	}

	slog.Info("Successfully processed message", "subject", msg.Subject)
}

// handleTopic1 processes messages from topic1
func (s *NatsSubscriberService) handleTopic1(tx *sql.Tx, data []byte) error {
	// Parse message data
	var messageData struct {
		// Define your message structure
		ID     string `json:"id"`
		Action string `json:"action"`
	}

	if err := json.Unmarshal(data, &messageData); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	// Process based on message content
	slog.Info("Processing topic1 message", "id", messageData.ID, "action", messageData.Action)

	// Implement your business logic here

	return nil
}

// handleTopic2 processes messages from topic2
func (s *NatsSubscriberService) handleTopic2(tx *sql.Tx, data []byte) error {
	// Similar to handleTopic1
	slog.Info("Processing topic2 message")

	// Implement your business logic here

	return nil
}

// Shutdown stops the subscriber service
func (s *NatsSubscriberService) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down NATS Subscriber Service")

	// Unsubscribe from all subscriptions
	for _, sub := range s.subscriptions {
		if sub != nil {
			if err := sub.Unsubscribe(); err != nil {
				slog.Error("Failed to unsubscribe", "error", err)
			}
		}
	}

	// Close connections
	if s.NatsClient != nil {
		s.NatsClient.Close()
	}

	return nil
}
