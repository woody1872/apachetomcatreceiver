package apachetomcat

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

var (
	typeStr = component.MustNewType("apachetomcat")
)

const (
	defaultInterval = 1 * time.Minute
)

// NewFactory creates a factory for the apachetomcat receiver
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, component.StabilityLevelAlpha))
}

// createDefaultConfig creates and returns the default configuration for the apachetomcat receiver
func createDefaultConfig() component.Config {
	return &Config{
		Interval: fmt.Sprint(defaultInterval),
	}
}

func createMetricsReceiver(_ context.Context, params receiver.CreateSettings, baseCfg component.Config, consumer consumer.Metrics) (receiver.Metrics, error) {
	logger := params.Logger
	apacheTomcatCfg := baseCfg.(*Config)

	metricRcvr := &apacheTomcatReceiver{
		logger:       logger,
		nextConsumer: consumer,
		config:       apacheTomcatCfg,
	}

	return metricRcvr, nil
}
