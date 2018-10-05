package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/satori/go.uuid"
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
func GetLoggers() ([]Logger, error) {
	results, err := load()
	if err != nil {
		return nil, err
	}

	loggers := make([]Logger, len(results))
	index := 0
	for _, value := range results {
		loggers[index] = value
		index++
	}

	return loggers, nil
}

// GetLogger with given ID
func GetLogger(loggerID string) (*Logger, error) {
	results, err := load()
	if err != nil {
		return nil, err
	}

	if item, ok := results[loggerID]; ok {
		return &item, nil
	}

	return nil, nil
}

var lock = &sync.Mutex{}

// AddLogger in configuration
func AddLogger(logger Logger) (string, error) {
	lock.Lock()
	defer lock.Unlock()
	logger.ID = uuid.Must(uuid.NewV4()).String()
	loggers, err := load()
	if err != nil {
		return "", err
	}

	loggers[logger.ID] = logger
	return logger.ID, save(loggers)
}

// UpdateLogger in configuration
func UpdateLogger(logger Logger) error {
	lock.Lock()
	defer lock.Unlock()
	loggers, err := load()
	if err != nil {
		return err
	}

	if _, ok := loggers[logger.ID]; ok {
		loggers[logger.ID] = logger
		return save(loggers)
	}

	return fmt.Errorf("unable to find logger with Id %s", logger.ID)
}

// DeleteLogger in configuration
func DeleteLogger(loggerID string) error {
	lock.Lock()
	defer lock.Unlock()
	loggers, err := load()
	if err != nil {
		return err
	}

	if _, ok := loggers[loggerID]; ok {
		delete(loggers, loggerID)
		return save(loggers)
	}

	return fmt.Errorf("unable to find logger with Id %s", loggerID)
}

// Setup the configuration
func Setup(file string) error {
	configFileName = file
	return nil
}

var configFileName string

func load() (map[string]Logger, error) {
	loggers := make(map[string]Logger)
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		return loggers, nil
	}

	fileData, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return loggers, err
	}

	err = json.Unmarshal(fileData, &loggers)
	return loggers, err
}

func save(config map[string]Logger) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFileName, configJSON, 0644)
}
