package loggers

import (
	"regexp"

	"github.com/cjburchell/go-uatu"
	"github.com/cjburchell/uatu/config"
)

// Logger item
type Logger struct {
	config.Logger
	re          *regexp.Regexp
	Destination Destination
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

// SetMaxLevel Sets the max level for a logger
func (l *Logger) SetMaxLevel(maxLevel int) {
	l.Levels = []int{}
	for i := 0; i <= maxLevel; i++ {
		l.Levels = append(l.Levels, i)
	}
}

// Check checks to see if the message should be logged
func (l Logger) Check(message log.Message) bool {

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
			var err error
			l.re, err = regexp.Compile(l.Pattern)
			if err != nil {
				l.Pattern = ""
			}
		}

		if !l.re.MatchString(message.Text) {
			return false
		}
	}

	return true
}

// UpdateDestination updates the destination
func (l *Logger) UpdateDestination() error {
	var err error
	l.Destination, err = destinations[l.DestinationType](l.DestinationConfig)
	if err != nil {
		return err
	}

	return l.Destination.Setup()
}

// Load the log file
func Load() ([]Logger, error) {
	result, err := config.GetLoggers()
	if err != nil {
		return nil, err
	}

	loggers := make([]Logger, len(result))
	for index, item := range result {
		loggers[index] = Logger{Logger: item}
		err = loggers[index].UpdateDestination()
		if err != nil {
			return nil, err
		}
	}

	return loggers, nil
}
