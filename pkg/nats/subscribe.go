package nats

import (
	"context"
	"errors"
	"time"

	nats "github.com/nats-io/nats.go"
)

// SubscribeOptions holds options for subscribing to subjects
type SubscribeOptions struct {
	Queue   string
	Timeout time.Duration
}

// DefaultSubscribeOptions returns default subscribe options
func DefaultSubscribeOptions() SubscribeOptions {
	return SubscribeOptions{
		Queue:   "",
		Timeout: time.Second * 30,
	}
}

// Subscribe subscribes to a subject
func (c *NATSClient) Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsConnected() {
		return nil, errors.New("not connected to NATS server")
	}

	return c.conn.Subscribe(subject, handler)
}

// SubscribeWithOptions subscribes to a subject with options
func (c *NATSClient) SubscribeWithOptions(subject string, handler nats.MsgHandler, opts SubscribeOptions) (*nats.Subscription, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsConnected() {
		return nil, errors.New("not connected to NATS server")
	}

	if opts.Queue != "" {
		return c.conn.QueueSubscribe(subject, opts.Queue, handler)
	}

	return c.conn.Subscribe(subject, handler)
}

// SubscribeSync subscribes to a subject and returns a subscription that can be used to manually fetch messages
func (c *NATSClient) SubscribeSync(subject string) (*nats.Subscription, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsConnected() {
		return nil, errors.New("not connected to NATS server")
	}

	return c.conn.SubscribeSync(subject)
}

// SubscribeWithContext subscribes to a subject with context for cancellation
func (c *NATSClient) SubscribeWithContext(ctx context.Context, subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsConnected() {
		return nil, errors.New("not connected to NATS server")
	}

	sub, err := c.conn.Subscribe(subject, handler)
	if err != nil {
		return nil, err
	}

	// Handle context cancellation
	go func() {
		<-ctx.Done()
		sub.Unsubscribe()
	}()

	return sub, nil
}

// Unsubscribe unsubscribes from a subscription
func (c *NATSClient) Unsubscribe(sub *nats.Subscription) error {
	if sub == nil {
		return errors.New("subscription is nil")
	}
	return sub.Unsubscribe()
}
