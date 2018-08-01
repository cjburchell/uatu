package logger

import (
	"fmt"
	"github.com/cjburchell/tools-go"
	"github.com/cjburchell/yasls-client-go"
	"github.com/robfig/cron"
	"gopkg.in/natefinch/lumberjack.v2"
)

func CreateFileDestination(maxAge, maxBackups, maxSize int, filename string) Destination {
	file := fileDestination{
		logger: lumberjack.Logger{
			MaxAge:     maxAge, //days
			MaxBackups: maxBackups,
			MaxSize:    maxSize, // megabytes
			Filename:   filename},
	}

	c := cron.New()
	c.AddFunc("@midnight", func() {
		fmt.Println("Resetting Logging file.")
		file.logger.Rotate()
	})
	c.Start()

	return file
}

type fileDestination struct {
	logger lumberjack.Logger
}

func (f fileDestination) PrintMessage(message log.LogMessage) {
	f.logger.Write([]byte(message.String() + "\n"))
}
