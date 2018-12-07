package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

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

type config struct {
	Loggers []Logger `json:"loggers"`
}

// GetLoggers configuration
func GetLoggers() ([]Logger, error) {
	return load()
}

// Setup the configuration
func Setup(file string) error {
	configFileName = file
	return nil
}

var configFileName string

func load() ([]Logger, error) {
	var err error
	if _, err = os.Stat(configFileName); os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}

	fileData, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var configData config
	err = json.Unmarshal(fileData, &configData)
	return configData.Loggers, errors.WithStack(err)
}
