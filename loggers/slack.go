package loggers

import (
	"github.com/cjburchell/yasls-client-go"
	"encoding/json"
	"github.com/nlopes/slack"
)

func createSlackDestination(data *json.RawMessage) Destination {
	var destination slackDestination
	if data == nil{
		return &destination
	}

	json.Unmarshal(*data, &destination)
	return &destination
}

type slackDestination struct {
	Token string `json:"token"`
	Channel string `json:"channel"`
	client *slack.Client
}

func (s slackDestination) PrintMessage(message log.LogMessage) {
	if s.client == nil{
		return
	}

	params := slack.PostMessageParameters{}
	s.client.PostMessage(s.Channel ,message.String(), params)
}

func (s *slackDestination) Stop()  {
	s.client = nil
}

func (s *slackDestination) Setup() error  {
	s.client = slack.New(s.Token)
	return nil
}

func init() {
	destinations["slack"] =  createSlackDestination
}