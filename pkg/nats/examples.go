package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	nats "github.com/nats-io/nats.go"
)

// ExampleMessage represents a sample message structure
type ExampleMessage struct {
	ID      string    `json:"id"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

// ExamplePublish demonstrates how to publish a message
func ExamplePublish(client *NATSClient) {
	// Create a sample message
	msg := ExampleMessage{
		ID:      "msg-123",
		Content: "Hello NATS",
		Time:    time.Now().UTC().Add(7 * time.Hour),
	}

	// Marshal the message to JSON
	data, err := json.Marshal(msg)
	if err != nil {
		slog.Error("Failed to marshal message", "error", err)
		return
	}

	// Publish the message
	err = client.Publish("example.subject", data)
	if err != nil {
		slog.Error("Failed to publish message", "error", err)
		return
	}

	slog.Info("Message published successfully", "subject", "example.subject")
}

// ExampleSubscribe demonstrates how to subscribe to a subject
func ExampleSubscribe(client *NATSClient) {
	// Create a message handler
	handler := func(msg *nats.Msg) {
		var receivedMsg ExampleMessage
		err := json.Unmarshal(msg.Data, &receivedMsg)
		if err != nil {
			slog.Error("Failed to unmarshal message", "error", err)
			return
		}

		// If the message has a reply subject, send a response
		if msg.Reply != "" {
			response := []byte("Message received")
			err = client.Reply(msg, response)
			if err != nil {
				slog.Error("Failed to send reply", "error", err)
			}
		}
	}

	// Subscribe to the subject
	_, err := client.Subscribe("example.subject", handler)
	if err != nil {
		slog.Error("Failed to subscribe", "error", err)
		return
	}

	slog.Info("Subscribed successfully", "subject", "example.subject")

	// To unsubscribe later:
	// client.Unsubscribe(sub)
}

// ExampleRequestReply demonstrates the request-reply pattern
func ExampleRequestReply(client *NATSClient) {
	// Set up a responder
	responder := func(msg *nats.Msg) {
		// Process the request
		// Send a reply
		response := fmt.Sprintf("Response to: %s", string(msg.Data))
		err := client.Reply(msg, []byte(response))
		if err != nil {
			slog.Error("Failed to send reply", "error", err)
		}
	}

	// Subscribe to handle requests
	_, err := client.Subscribe("example.request", responder)
	if err != nil {
		slog.Error("Failed to subscribe for requests", "error", err)
		return
	}

	// Send a request and wait for reply
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reply, err := client.RequestWithContext(ctx, "example.request", []byte("Hello, need a response"))
	if err != nil {
		slog.Error("Request failed", "error", err)
		return
	}

	slog.Info("Received reply", "data", string(reply.Data))
}
