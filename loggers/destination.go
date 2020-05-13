package loggers

import (
	"encoding/json"

	uatu "github.com/cjburchell/uatu-go"
)

// Destination for a logger
type Destination interface {
	PrintMessage(message uatu.Message) error
	Setup() error
	Stop()
}

var destinations = map[string]func(*json.RawMessage) (Destination, error){}
