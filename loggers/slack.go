package loggers

import (
	"encoding/json"

	"github.com/cjburchell/yasls-client-go"
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
	Token   string `json:"token"`
	Channel string `json:"channel"`
	client  *slack.Client
}

func (s slackDestination) PrintMessage(message log.LogMessage) error {
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
	s.client = slack.New(s.Token)
	return nil
}

func init() {
	destinations["slack"] = createSlackDestination
}
