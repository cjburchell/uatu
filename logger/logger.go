package logger

import "github.com/cjburchell/yasls-client-go"

type Logger interface {
	PrintMessage(message log.LogMessage)
}

type loggerBase struct {
	enabled  bool
	minLevel int
}
