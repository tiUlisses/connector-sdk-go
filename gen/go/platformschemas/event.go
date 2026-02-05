package platformschemas

import "encoding/json"

// Event is a minimal schema payload representing serialized platform events.
type Event struct {
	Type      string          `json:"type"`
	Timestamp int64           `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

// SerializeEvent serializes an Event as JSON bytes.
func SerializeEvent(event Event) ([]byte, error) {
	return json.Marshal(event)
}
