package loggers

import (
	"fmt"
	"github.com/cjburchell/yasls-client-go"
	"encoding/json"
)

func createConsoleDestination(_ *json.RawMessage) Destination {
	return consoleDestination{}
}

type consoleDestination struct {
}

func (consoleDestination) PrintMessage(message log.LogMessage) {
	fmt.Println(message.String())
}

func (consoleDestination) Stop()  {
}

func (consoleDestination) Setup() error  {
	return nil
}

func init() {
	destinations["console"] =  createConsoleDestination
}
