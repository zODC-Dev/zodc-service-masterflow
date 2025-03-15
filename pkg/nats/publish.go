package nats

import (
	"context"
	"errors"
	"time"

	nats "github.com/nats-io/nats.go"
)

// PublishOptions holds options for publishing messages
type PublishOptions struct {
	Headers nats.Header
	Timeout time.Duration
}

// DefaultPublishOptions returns default publish options
func DefaultPublishOptions() PublishOptions {
	return PublishOptions{
		Headers: nats.Header{},
		Timeout: time.Second * 5,
	}
}

// Publish publishes a message to the specified subject
func (c *NATSClient) Publish(subject string, data []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsConnected() {
		return errors.New("not connected to NATS server")
	}

	return c.conn.Publish(subject, data)
}

// PublishWithOptions publishes a message with the specified options
func (c *NATSClient) PublishWithOptions(subject string, data []byte, opts PublishOptions) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsConnected() {
		return errors.New("not connected to NATS server")
	}

	msg := &nats.Msg{
		Subject: subject,
		Data:    data,
		Header:  opts.Headers,
	}

	return c.conn.PublishMsg(msg)
}

// PublishAsync publishes a message asynchronously
func (c *NATSClient) PublishAsync(subject string, data []byte) (nats.PubAckFuture, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsConnected() {
		return nil, errors.New("not connected to NATS server")
	}

	if c.js == nil {
		return nil, errors.New("JetStream not enabled")
	}

	return c.js.PublishAsync(subject, data)
}

// PublishWithContext publishes a message with context for timeout/cancellation
func (c *NATSClient) PublishWithContext(ctx context.Context, subject string, data []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsConnected() {
		return errors.New("not connected to NATS server")
	}

	// Create a channel to signal completion
	done := make(chan error, 1)

	go func() {
		done <- c.conn.Publish(subject, data)
	}()

	// Wait for either context cancellation or publish completion
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
