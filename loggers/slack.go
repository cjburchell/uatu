package loggers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/cjburchell/go-uatu"
)

func createSlackDestination(data *json.RawMessage) (Destination, error) {
	var destination slackDestination
	if data == nil {
		return &destination, nil
	}

	err := json.Unmarshal(*data, &destination)
	return &destination, err
}

type slackDestination struct {
	Destination string `json:"destination"`
}

func (s slackDestination) PrintMessage(message log.Message) error {
	return s.sendMessage(message.String())
}

func (s slackDestination) sendMessage(message string) error {
	jsonValue, err := json.Marshal(struct {
		Text string `json:"text"`
	}{
		Text: message,
	})

	if err != nil {
		return errors.WithStack(err)
	}

	resp, err := http.Post(s.Destination, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return errors.WithStack(err)
	}

	if resp.StatusCode != http.StatusOK {
		return errors.WithStack(fmt.Errorf("http request to slack %s error: %d", s.Destination, resp.StatusCode))
	}

	return nil
}

func (s *slackDestination) Stop() {
	s.sendMessage("Stop Logging")
}

func (s *slackDestination) Setup() error {
	return s.sendMessage("Start Logging")
}

func init() {
	destinations["slack"] = createSlackDestination
}
