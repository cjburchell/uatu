package loggers

import (
	"encoding/json"

	"github.com/cjburchell/tools-go/env"

	"github.com/cjburchell/go-uatu"
	"github.com/nlopes/slack"
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
	Channel string `json:"channel"`
	client  *slack.Client
}

func (s slackDestination) PrintMessage(message log.Message) error {
	if s.client == nil {
		return nil
	}

	params := slack.PostMessageParameters{}
	_, _, err := s.client.PostMessage(s.Channel, message.String(), params)
	return err
}

func (s *slackDestination) Stop() {
	s.client = nil
}

func (s *slackDestination) Setup() error {
	token := env.Get("SLACK_TOKEN", "")
	s.client = slack.New(token)
	return nil
}

func init() {
	destinations["slack"] = createSlackDestination
}
