package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/example/connector-sdk-go/pkg/connector"
)

type MinimalConnector struct{}

func (m *MinimalConnector) Metadata() connector.Metadata {
	return connector.Metadata{Name: "minimal", Version: "0.1.0", Capabilities: []string{"commands"}}
}

func (m *MinimalConnector) ValidateConfig(raw []byte) error {
	if len(raw) == 0 {
		return errors.New("empty config")
	}
	var temp map[string]any
	return json.Unmarshal(raw, &temp)
}

func (m *MinimalConnector) Discover(context.Context) ([]connector.DeviceDiscovered, error) {
	return []connector.DeviceDiscovered{{ID: "device-001", Name: "Example Device"}}, nil
}

func (m *MinimalConnector) Start(context.Context, string, []byte) error { return nil }

func (m *MinimalConnector) Stop(context.Context, string) error { return nil }

func (m *MinimalConnector) ExecuteCommand(_ context.Context, deviceID, commandName string, payload []byte) ([]byte, error) {
	response := map[string]any{
		"deviceId":     deviceID,
		"command":      commandName,
		"payloadBytes": len(payload),
		"processedAt":  time.Now().UTC().Format(time.RFC3339),
	}
	return json.Marshal(response)
}

func main() {
	var c connector.Connector = &MinimalConnector{}
	resp, err := c.ExecuteCommand(context.Background(), "device-001", "ping", []byte(`{"value":1}`))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(resp))
}
