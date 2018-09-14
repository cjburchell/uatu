package loggers

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/cjburchell/yasls-client-go"
	"github.com/robfig/cron"
)

func createS3Destination(data *json.RawMessage) (Destination, error) {
	var destination s3Destination
	if data == nil {
		return &destination, nil
	}

	err := json.Unmarshal(*data, &destination)
	return &destination, err
}

type s3Destination struct {
	Bucket     string `json:"bucket"`
	Region     string `json:"region"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
	MaxSize    int    `json:"max_size"`
	Filename   string `json:"filename"`
	session    *session.Session
	cron       *cron.Cron
}

func (s s3Destination) PrintMessage(message log.LogMessage) error {
	return nil
}

func (s3Destination) Stop() {
}

func (s *s3Destination) Setup() error {
	var err error
	s.session, err = session.NewSession(&aws.Config{Region: aws.String(s.Region)})
	if err != nil {
		return err
	}

	s.cron = cron.New()
	err = s.cron.AddFunc("@midnight", func() {
		fmt.Println("Resetting Logging file.")
	})

	if err != nil {
		return err
	}

	s.cron.Start()

	return nil
}

func init() {
	destinations["s3"] = createS3Destination
}
