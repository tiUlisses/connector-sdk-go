package connector

import (
	"context"
)

// Metadata contains static connector identification information.
type Metadata struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Capabilities []string `json:"capabilities,omitempty"`
}

// DeviceDiscovered represents a device found during discovery.
type DeviceDiscovered struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Labels   map[string]string `json:"labels,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Connector defines the required operations for an SDK connector implementation.
type Connector interface {
	Metadata() Metadata
	ValidateConfig(raw []byte) error
	Discover(ctx context.Context) ([]DeviceDiscovered, error)
	Start(ctx context.Context, deviceID string, config []byte) error
	Stop(ctx context.Context, deviceID string) error
	ExecuteCommand(ctx context.Context, deviceID, commandName string, payload []byte) ([]byte, error)
}
