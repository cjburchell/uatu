package logger

import (
	"fmt"
	"github.com/cjburchell/tools-go"
	"github.com/cjburchell/yasls-client-go"
)

func CreateConsole() Logger {
	return console{
		loggerBase: loggerBase{
			minLevel: tools.GetEnvInt("LOG_CONSOLE_LEVEL", log.INFO.Severity),
			enabled:  tools.GetEnvBool("LOG_CONSOLE", false),
		},
	}
}

type console struct {
	loggerBase
}

func (c console) PrintMessage(message log.LogMessage) {
	if message.Level.Severity >= c.minLevel && c.enabled {
		fmt.Println(message.String())
	}
}
