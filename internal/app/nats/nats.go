package nats

import (
	"log/slog"
	"sync"

	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/configs"
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
