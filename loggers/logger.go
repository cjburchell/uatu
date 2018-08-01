package loggers

import (
	"github.com/cjburchell/yasls-client-go"
	"regexp"
	"encoding/json"
	"io/ioutil"
)

type Logger struct {
	Id              string           `json:"id"`
	Description     string           `json:"description, omitempty"`
	Levels          []int            `json:"levels, omitempty"`
	Sources         []string         `json:"sources, omitempty"`
	Pattern         string           `json:"pattern, omitempty"`
	DestinationType string           `json:"destination_type"`
	DestConfig      *json.RawMessage `json:"destination, omitempty"`
	re              *regexp.Regexp   `json:"-"`
	Destination     Destination      `json:"-"`
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
	l.Destination = destinations[l.DestinationType](l.DestConfig)
	l.Destination.Setup()
}

func (l *Logger) UpdateDestinationConfig() {
	data, _:= json.Marshal(l.Destination)
	config := json.RawMessage(data)
	l.DestConfig = &config
}

func Load(file string) ([]Logger, error) {
	var loggers []Logger
	fileData, err:= ioutil.ReadFile(file)
	if err != nil{
		return loggers, err
	}

	err = json.Unmarshal(fileData, &loggers)
	return loggers, err
}


