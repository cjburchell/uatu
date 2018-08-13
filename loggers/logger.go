package loggers

import (
	"github.com/cjburchell/yasls-client-go"
	"regexp"
	"github.com/cjburchell/yasls/config"
)

type Logger struct {
	config.Logger
	re              *regexp.Regexp
	Destination     Destination
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (l *Logger) SetMaxLevel(maxLevel int) {
	l.Levels = []int{}
	for i := 0; i <= maxLevel; i++ {
		l.Levels = append(l.Levels, i)
	}
}

func (l Logger) Check(message log.LogMessage) bool {

	if len(l.Levels) != 0 {
		if !intInSlice(message.Level.Severity, l.Levels) {
			return false
		}
	}

	if len(l.Sources) != 0 {
		if !stringInSlice(message.ServiceName, l.Sources) {
			return false
		}
	}

	if l.Pattern != "" {
		if l.re == nil {
			l.re, _ = regexp.Compile(l.Pattern)
		}

		if !l.re.MatchString(message.Text) {
			return false
		}
	}

	return true
}

func (l* Logger) UpdateDestination()  {
	l.Destination = destinations[l.DestinationType](l.DestinationConfig)
	l.Destination.Setup()
}

func Load() ([]Logger, error)  {
	result, err := config.GetLoggers()
	if err != nil{
		return nil, err
	}

	loggers := make([]Logger, len(result))
	for index, item := range result {
		loggers[index] = Logger{Logger: item }
		loggers[index].UpdateDestination()
	}

	return loggers, nil
}
