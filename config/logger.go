package config

import (
	"encoding/json"
	"io/ioutil"
	"github.com/satori/go.uuid"
	"fmt"
	"sync"
)

type Logger struct {
	Id                string           `json:"id"`
	Description       string           `json:"description, omitempty"`
	Levels            []int            `json:"levels, omitempty"`
	Sources           []string         `json:"sources, omitempty"`
	Pattern           string           `json:"pattern, omitempty"`
	DestinationType   string           `json:"destination_type"`
	DestinationConfig *json.RawMessage `json:"destination, omitempty"`
}

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

func GetLogger(loggerId string) (*Logger, error)  {
	results, err := load()
	if err != nil{
		return nil, err
	}

	if item, ok := results[loggerId]; ok{
		return &item, nil
	}

	return nil, nil
}

var lock = &sync.Mutex{}

func AddLogger(logger Logger) (string, error)  {
	lock.Lock()
	defer lock.Unlock()
	logger.Id = uuid.Must(uuid.NewV4()).String()
	loggers, err := load()
	if err != nil{
		return "", err
	}

	loggers[logger.Id] = logger
	return logger.Id, save(loggers)
}

func UpdateLogger(logger Logger) error  {
	lock.Lock()
	defer lock.Unlock()
	loggers, err := load()
	if err != nil{
		return err
	}

	if _, ok := loggers[logger.Id]; ok{
		loggers[logger.Id] = logger
		return save(loggers)
	}

	return fmt.Errorf("unable to find logger with Id %s", logger.Id)
}

func DeleteLogger(loggerId string) error {
	lock.Lock()
	defer lock.Unlock()
	loggers, err := load()
	if err != nil {
		return err
	}

	if _, ok := loggers[loggerId]; ok {
		delete(loggers, loggerId)
		return save(loggers)
	}

	return fmt.Errorf("unable to find logger with Id %s", loggerId)
}

func Setup(file string) error  {
	configFileName = file
	return nil
}

var configFileName string

func load() (map[string]Logger, error) {
	loggers := make(map[string]Logger)
	fileData, err:= ioutil.ReadFile(configFileName)
	if err != nil{
		return loggers, err
	}

	err = json.Unmarshal(fileData, &loggers)
	return loggers, err
}

func save(config map[string]Logger) error {
	configJson, _ := json.Marshal(config)
	return ioutil.WriteFile(configFileName, configJson, 0644)
}