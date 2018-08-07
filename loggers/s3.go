package loggers

import (
	"encoding/json"
	"github.com/cjburchell/yasls-client-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/robfig/cron"
	"fmt"
)

func createS3Destination(data *json.RawMessage) Destination {
	var destination s3Destination
	if data == nil{
		return &destination
	}

	json.Unmarshal(*data, &destination)
	return &destination
}

type s3Destination struct {
	Bucket     string `json:"bucket"`
	Region     string `json:"region"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
	MaxSize    int    `json:"max_size"`
	Filename   string `json:"filename"`
	session    *session.Session
}

func (s s3Destination) PrintMessage(message log.LogMessage) {
}

func (s3Destination) Stop()  {
}

func (s *s3Destination) Setup() error {
	var err error
	s.session, err = session.NewSession(&aws.Config{Region: aws.String(s.Region)})
	if err != nil {
		return err
	}

	f.cron = cron.New()
	f.cron.AddFunc("@midnight", func() {
		fmt.Println("Resetting Logging file.")
		f.logger.Rotate()
	})
	f.cron.Start()

	return nil
}

func init() {
	destinations["s3"] =  createS3Destination
}
