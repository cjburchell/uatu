package loggers

import (
	"encoding/json"
	"fmt"

	"github.com/cjburchell/yasls-client-go"
)

func createConsoleDestination(_ *json.RawMessage) (Destination, error) {
	return consoleDestination{}, nil
}

type consoleDestination struct {
}

func (consoleDestination) PrintMessage(message log.LogMessage) error {
	_, err := fmt.Println(message.String())
	return err
}

func (consoleDestination) Stop() {
}

func (consoleDestination) Setup() error {
	return nil
}

func init() {
	destinations["console"] = createConsoleDestination
}
