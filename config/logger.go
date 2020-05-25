package config

import (
	"encoding/json"
	"github.com/cjburchell/settings-go"
	"github.com/pkg/errors"
)

// Logger configuration
type Logger struct {
	ID                string           `json:"id"`
	Description       string           `json:"description"`
	Levels            []int            `json:"levels"`
	Sources           []string         `json:"sources"`
	Pattern           string           `json:"pattern"`
	DestinationType   string           `json:"destination_type"`
	DestinationConfig *json.RawMessage `json:"destination"`
}

// GetLoggers configuration
func GetLoggers(settings settings.ISettings) ([]Logger, error) {
	return load(settings)
}

func load(settings settings.ISettings) ([]Logger, error) {
	var loggers []Logger
	err := settings.GetObject("loggers", &loggers)
	return loggers, errors.WithStack(err)
}
