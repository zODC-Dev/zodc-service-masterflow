package nats

import (
	"log/slog"
	"sync"
	"time"

	nats "github.com/nats-io/nats.go"
)

// NATSClient represents a NATS client with connection management
type NATSClient struct {
	conn      *nats.Conn
	js        nats.JetStreamContext
	mu        sync.RWMutex
	config    Config
	reconnect bool
}

// Config holds the configuration for NATS client
type Config struct {
	URL            string
	MaxReconnects  int
	ReconnectWait  time.Duration
	ConnectTimeout time.Duration
	UseJetStream   bool
}

// DefaultConfig returns a default configuration for NATS client
func DefaultConfig() Config {
	return Config{
		URL:            nats.DefaultURL,
		MaxReconnects:  10,
		ReconnectWait:  time.Second * 2,
		ConnectTimeout: time.Second * 5,
		UseJetStream:   false,
	}
}

// NewNATSClient creates a new NATS client with the given configuration
func NewNATSClient(config Config) (*NATSClient, error) {
	client := &NATSClient{
		config:    config,
		reconnect: true,
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	return client, nil
}

// connect establishes a connection to the NATS server
func (c *NATSClient) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	opts := []nats.Option{
		nats.MaxReconnects(c.config.MaxReconnects),
		nats.ReconnectWait(c.config.ReconnectWait),
		nats.Timeout(c.config.ConnectTimeout),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			slog.Error("NATS disconnected", "error", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			slog.Info("NATS reconnected", "url", nc.ConnectedUrl())
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			slog.Error("NATS error", "error", err, "subject", sub.Subject)
		}),
	}

	var err error
	c.conn, err = nats.Connect(c.config.URL, opts...)
	if err != nil {
		return err
	}

	if c.config.UseJetStream {
		c.js, err = c.conn.JetStream()
		if err != nil {
			c.conn.Close()
			return err
		}
	}

	slog.Info("Connected to NATS server", "url", c.conn.ConnectedUrl())
	return nil
}

// Close closes the NATS connection
func (c *NATSClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.reconnect = false
	if c.conn != nil {
		c.conn.Close()
	}
	slog.Info("NATS connection closed")
}

// GetConn returns the underlying NATS connection
func (c *NATSClient) GetConn() *nats.Conn {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn
}

// GetJetStream returns the JetStream context
func (c *NATSClient) GetJetStream() nats.JetStreamContext {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.js
}

// IsConnected returns true if the client is connected to the NATS server
func (c *NATSClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn != nil && c.conn.IsConnected()
}
