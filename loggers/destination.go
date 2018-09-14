package loggers

import (
	"encoding/json"

	"github.com/cjburchell/yasls-client-go"
)

// Destination for a logger
type Destination interface {
	PrintMessage(message log.LogMessage) error
	Setup() error
	Stop()
}

var destinations = map[string]func(*json.RawMessage) (Destination, error){}
