package kafka

import (
	"context"
	"errors"
	"testing"
	"time"

	kg "github.com/segmentio/kafka-go"

	"github.com/example/connector-sdk-go/gen/go/platformschemas"
)

type fakeWriter struct {
	failures int
	calls    int
	msg      kg.Message
}

func (f *fakeWriter) WriteMessages(_ context.Context, msgs ...kg.Message) error {
	f.calls++
	f.msg = msgs[0]
	if f.calls <= f.failures {
		return errors.New("temporary")
	}
	return nil
}

func (f *fakeWriter) Close() error { return nil }

func TestPublishRetriesAndHeaders(t *testing.T) {
	w := &fakeWriter{failures: 1}
	p := NewPublisherWithWriter(w, 3, time.Millisecond)

	err := p.Publish(context.Background(), MessageContext{
		TenantID:      "t1",
		SiteID:        "s1",
		DeviceID:      "d1",
		SchemaVersion: "v1",
	}, platformschemas.Event{Type: "x", Timestamp: 1, Payload: []byte(`{"a":1}`)})
	if err != nil {
		t.Fatalf("publish should succeed after retry: %v", err)
	}
	if w.calls != 2 {
		t.Fatalf("expected 2 calls, got %d", w.calls)
	}
	if len(w.msg.Headers) != 4 {
		t.Fatalf("expected 4 headers, got %d", len(w.msg.Headers))
	}
}
