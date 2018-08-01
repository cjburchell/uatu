package logger

import (
	"fmt"
	"github.com/cjburchell/tools-go"
	"github.com/cjburchell/yasls-client-go"
)

func CreateConsoleDestination() Destination {
	return consoleDestination{}
}

type consoleDestination struct {
}

func (consoleDestination) PrintMessage(message log.LogMessage) {
	fmt.Println(message.String())
}
