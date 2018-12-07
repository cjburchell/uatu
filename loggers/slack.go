package loggers

import (
	"encoding/json"

	"github.com/cjburchell/uatu/settings"

	"github.com/pkg/errors"

	"github.com/bluele/slack"
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
	Channel string `json:"channel"`
	client  *slack.Slack
}

func (s slackDestination) PrintMessage(message log.Message) error {
	if s.client == nil {
		return nil
	}

	err := s.client.ChatPostMessage(s.Channel, message.String(), nil)
	return errors.Wrapf(errors.WithStack(err), "Unable to post slack message to %s", s.Channel)
}

func (s *slackDestination) Stop() {
	if s.client == nil {
		return
	}

	s.client.ChatPostMessage(s.Channel, "Stop Logging", nil)
	s.client = nil
}

func (s *slackDestination) Setup() error {
	s.client = slack.New(settings.SlackToken)
	err := s.client.ChatPostMessage(s.Channel, "Start Logging", nil)
	return errors.Wrapf(errors.WithStack(err), "Unable to post slack message to %s", s.Channel)
}

func init() {
	destinations["slack"] = createSlackDestination
}
