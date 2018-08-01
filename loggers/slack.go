package loggers

import (
	"fmt"
	"github.com/cjburchell/yasls-client-go"
	"encoding/json"
	"github.com/nlopes/slack"
)

func createSlackDestination(_ *json.RawMessage) Destination {
	return slackDestination{}
}

type slackDestination struct {
	client *slack.Client
}

func (s slackDestination) PrintMessage(message log.LogMessage) {
	s.client.PostMessage("test" ,message.String(), )
}

func (slackDestination) Stop()  {
}

func (s *slackDestination) Setup()  {
	s.client = slack.New("token")
}