package loggers

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	uatu "github.com/cjburchell/uatu-go"
)

func createConsoleDestination(_ *json.RawMessage) (Destination, error) {
	return consoleDestination{}, nil
}

type consoleDestination struct {
}

func (consoleDestination) PrintMessage(message uatu.Message) error {
	_, err := fmt.Println(message.String())
	return errors.WithStack(err)
}

func (consoleDestination) Stop() {
}

func (consoleDestination) Setup() error {
	_, err := fmt.Println("Start Console Log")
	return errors.WithStack(err)
}

func init() {
	destinations["console"] = createConsoleDestination
}
