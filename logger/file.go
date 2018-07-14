package logger

import (
	"fmt"
	"github.com/cjburchell/tools-go"
	"github.com/cjburchell/yasls-client-go"
	"github.com/robfig/cron"
	"gopkg.in/natefinch/lumberjack.v2"
)

func CreateFile() Logger {
	file := file{
		logger: lumberjack.Logger{
			MaxAge:     tools.GetEnvInt("LOG_FILE_MAX_AGE", 1), //days
			MaxBackups: tools.GetEnvInt("LOG_FILE_MAX_BACKUPS", 20),
			MaxSize:    tools.GetEnvInt("LOG_FILE_MAX_SIZE", 100) * 1024 * 1024., // megabytes
			Filename:   tools.GetEnv("LOG_FILE_PATH", "/logs/server.log")},

		loggerBase: loggerBase{
			minLevel: tools.GetEnvInt("LOG_FILE_LEVEL", log.INFO.Severity),
			enabled:  tools.GetEnvBool("LOG_FILE", false),
		},
	}

	c := cron.New()
	c.AddFunc("@midnight", func() {
		fmt.Println("Resetting Logging file.")
		file.logger.Rotate()
	})
	c.Start()

	return file
}

type file struct {
	loggerBase
	logger lumberjack.Logger
}

func (f file) PrintMessage(message log.LogMessage) {
	if message.Level.Severity >= f.minLevel && f.enabled {
		f.logger.Write([]byte(message.String() + "\n"))
	}
}
