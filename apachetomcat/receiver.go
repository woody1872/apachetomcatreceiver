package apachetomcat

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

type apacheTomcatReceiver struct {
	host         component.Host
	cancel       context.CancelFunc
	logger       *zap.Logger
	nextConsumer consumer.Metrics
	config       *Config
}

func (atr *apacheTomcatReceiver) Start(ctx context.Context, host component.Host) error {
	atr.host = host
	ctx = context.Background()
	ctx, atr.cancel = context.WithCancel(ctx)
	client := newDefaultApacheTomcatClient(atr.config)

	interval, _ := time.ParseDuration(atr.config.Interval)
	go func(client *apacheTomcatClient) {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				tomcatStatus, err := client.getTomcatStatus()
				if err != nil {
					atr.logger.Error(fmt.Sprintf("%s: %v", "error reading Tomcat server status response body", err))
				}
				atr.logger.Debug(fmt.Sprintf("STATUS: %+v", tomcatStatus))

			case <-ctx.Done():
				return
			}
		}
	}(client)

	return nil
}

func (atr *apacheTomcatReceiver) Shutdown(ctx context.Context) error {
	atr.cancel()
	return nil
}
