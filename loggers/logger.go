package loggers

import (
	"log"
	"regexp"

	uatu "github.com/cjburchell/uatu-go"
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
func (l Logger) Check(message uatu.Message) bool {

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
func Load(file string) ([]Logger, error) {
	result, err := config.GetLoggers(file)
	if err != nil {
		return nil, err
	}

	var loggers []Logger
	for _, item := range result {
		logger := Logger{Logger: item}
		err = logger.UpdateDestination()
		if err != nil {
			log.Printf("Unable to setup logger %s, %s", logger.DestinationType, err.Error())

		} else {
			loggers = append(loggers, logger)
		}
	}

	return loggers, nil
}
