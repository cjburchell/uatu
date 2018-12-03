package loggers

import (
	"encoding/json"

	"github.com/cjburchell/go-uatu"
	"gopkg.in/natefinch/lumberjack.v2"
)

func createFileDestination(data *json.RawMessage) (Destination, error) {
	var file fileDestination
	if data == nil {
		return &file, nil
	}

	err := json.Unmarshal(*data, &file)
	return &file, err
}

type fileDestination struct {
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
	MaxSize    int    `json:"max_size"`
	Filename   string `json:"filename"`
	logger     *lumberjack.Logger
}

func (f fileDestination) PrintMessage(message log.Message) error {
	if f.logger != nil {
		_, err := f.logger.Write([]byte(message.String() + "\n"))
		return err
	}

	return nil
}

func (f *fileDestination) Setup() error {
	f.logger = &lumberjack.Logger{
		MaxAge:     f.MaxAge, //days
		MaxBackups: f.MaxBackups,
		MaxSize:    f.MaxSize, // megabytes
		Filename:   f.Filename}
	return nil
}

func (f *fileDestination) Stop() {
}

func init() {
	destinations["file"] = createFileDestination
}
