package nats

import (
	"context"
	"errors"
	"time"

	nats "github.com/nats-io/nats.go"
)

// RequestOptions holds options for request-reply operations
type RequestOptions struct {
	Headers nats.Header
	Timeout time.Duration
}

// DefaultRequestOptions returns default request options
func DefaultRequestOptions() RequestOptions {
	return RequestOptions{
		Headers: nats.Header{},
		Timeout: time.Second * 5,
	}
}

// Request sends a request and waits for a reply
func (c *NATSClient) Request(subject string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsConnected() {
		return nil, errors.New("not connected to NATS server")
	}

	return c.conn.Request(subject, data, timeout)
}

// RequestWithOptions sends a request with options and waits for a reply
func (c *NATSClient) RequestWithOptions(subject string, data []byte, opts RequestOptions) (*nats.Msg, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsConnected() {
		return nil, errors.New("not connected to NATS server")
	}

	msg := &nats.Msg{
		Subject: subject,
		Data:    data,
		Header:  opts.Headers,
	}

	return c.conn.RequestMsg(msg, opts.Timeout)
}

// RequestWithContext sends a request with context for timeout/cancellation
func (c *NATSClient) RequestWithContext(ctx context.Context, subject string, data []byte) (*nats.Msg, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsConnected() {
		return nil, errors.New("not connected to NATS server")
	}

	return c.conn.RequestWithContext(ctx, subject, data)
}

// Reply sends a reply to a message
func (c *NATSClient) Reply(msg *nats.Msg, data []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsConnected() {
		return errors.New("not connected to NATS server")
	}

	return c.conn.Publish(msg.Reply, data)
}

// ReplyWithHeaders sends a reply with headers to a message
func (c *NATSClient) ReplyWithHeaders(msg *nats.Msg, data []byte, headers nats.Header) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsConnected() {
		return errors.New("not connected to NATS server")
	}

	replyMsg := &nats.Msg{
		Subject: msg.Reply,
		Data:    data,
		Header:  headers,
	}

	return c.conn.PublishMsg(replyMsg)
}
