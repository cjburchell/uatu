package loggers

import (
	"github.com/cjburchell/yasls-client-go"
	"encoding/json"
)

type Destination interface {
	PrintMessage(message log.LogMessage)
	Setup()
	Stop()
}

var destinations = map[string]func(*json.RawMessage)Destination{
	"console": createConsoleDestination,
	"file": createFileDestination,
	"slack": createSlackDestination,
}
