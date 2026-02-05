package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/example/connector-sdk-go/gen/go/platformschemas"
)

// HeaderKey values used when publishing events.
const (
	HeaderTenantID      = "tenantId"
	HeaderSiteID        = "siteId"
	HeaderDeviceID      = "deviceId"
	HeaderSchemaVersion = "schemaVersion"
)

// MessageContext defines metadata headers sent with each event.
type MessageContext struct {
	TenantID      string
	SiteID        string
	DeviceID      string
	SchemaVersion string
}

type writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

// Publisher sends connector events to Kafka with retry/backoff.
type Publisher struct {
	w          writer
	maxRetries int
	backoff    time.Duration
}

// Config configures Kafka publishing behavior.
type Config struct {
	Brokers    []string
	Topic      string
	MaxRetries int
	Backoff    time.Duration
}

// NewPublisher constructs a kafka-backed Publisher.
func NewPublisher(cfg Config) (*Publisher, error) {
	if len(cfg.Brokers) == 0 {
		return nil, fmt.Errorf("brokers are required")
	}
	if cfg.Topic == "" {
		return nil, fmt.Errorf("topic is required")
	}

	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 3
	}
	if cfg.Backoff <= 0 {
		cfg.Backoff = 500 * time.Millisecond
	}

	kw := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &Publisher{w: kw, maxRetries: cfg.MaxRetries, backoff: cfg.Backoff}, nil
}

// NewPublisherWithWriter is intended for tests.
func NewPublisherWithWriter(w writer, maxRetries int, backoff time.Duration) *Publisher {
	if maxRetries <= 0 {
		maxRetries = 3
	}
	if backoff <= 0 {
		backoff = 500 * time.Millisecond
	}
	return &Publisher{w: w, maxRetries: maxRetries, backoff: backoff}
}

// Publish serializes and publishes an event with mandatory headers.
func (p *Publisher) Publish(ctx context.Context, msgCtx MessageContext, event platformschemas.Event) error {
	payload, err := platformschemas.SerializeEvent(event)
	if err != nil {
		return fmt.Errorf("serialize event: %w", err)
	}

	msg := kafka.Message{
		Value: payload,
		Headers: []kafka.Header{
			{Key: HeaderTenantID, Value: []byte(msgCtx.TenantID)},
			{Key: HeaderSiteID, Value: []byte(msgCtx.SiteID)},
			{Key: HeaderDeviceID, Value: []byte(msgCtx.DeviceID)},
			{Key: HeaderSchemaVersion, Value: []byte(msgCtx.SchemaVersion)},
		},
	}

	var lastErr error
	for attempt := 0; attempt < p.maxRetries; attempt++ {
		if err := p.w.WriteMessages(ctx, msg); err == nil {
			return nil
		} else {
			lastErr = err
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("publish canceled: %w", ctx.Err())
		case <-time.After(p.backoff * time.Duration(attempt+1)):
		}
	}

	return fmt.Errorf("publish failed after retries: %w", lastErr)
}

// Close closes underlying Kafka writer resources.
func (p *Publisher) Close() error {
	if p.w == nil {
		return errors.New("publisher writer is nil")
	}
	return p.w.Close()
}
