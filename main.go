package main

import (
	"encoding/json"
	"fmt"
	"github.com/cjburchell/tools-go"
	"github.com/cjburchell/yasls-client-go"
	"github.com/cjburchell/yasls/logger"
	"github.com/cjburchell/yasls/web"
	"github.com/nats-io/go-nats"
)

func main() {
	natsUrl := tools.GetEnv("NATS_URL", "tcp://nats:4222")

	maxAge := tools.GetEnvInt("LOG_FILE_MAX_AGE", 1) //days
	maxBackups := tools.GetEnvInt("LOG_FILE_MAX_BACKUPS", 20)
	maxSize := tools.GetEnvInt("LOG_FILE_MAX_SIZE", 100) * 1024 * 1024. // megabytes
	filename := tools.GetEnv("LOG_FILE_PATH", "/logs/server.log")

	minLevelFile := tools.GetEnvInt("LOG_FILE_LEVEL", log.INFO.Severity)
	enabledFile := tools.GetEnvBool("LOG_FILE", false)

	if enabledFile {
		consoleLogger := logger.Logger{Destination: logger.CreateFileDestination(maxAge, maxBackups, maxSize, filename)}
		consoleLogger.SetMaxLevel(minLevelFile)
		loggers = append(loggers, consoleLogger)
	}

	minLevelConsole := tools.GetEnvInt("LOG_CONSOLE_LEVEL", log.INFO.Severity)
	enabledConsole := tools.GetEnvBool("LOG_CONSOLE", false)
	if enabledConsole {
		consoleLogger := logger.Logger{Destination: logger.CreateConsoleDestination()}
		consoleLogger.SetMaxLevel(minLevelConsole)
		loggers = append(loggers, consoleLogger)
	}

	err := setupNats(natsUrl)
	if err != nil {
		fmt.Printf("Can't connect to nats: %v\n", err)
		return
	}
	web.StartHttp()
}

var loggers []logger.Logger

var natsConn *nats.Conn

func setupNats(natsUrl string) error {
	fmt.Println("Connecting to Nats " + natsUrl)
	var err error
	natsConn, err = nats.Connect(natsUrl)
	if err != nil {
		return err
	}

	_, err = natsConn.Subscribe("logs", handleMessage)
	natsConn.Flush()
	return err
}

func handleMessage(message *nats.Msg) {
	data := message.Data
	logMessage := log.LogMessage{}
	if err := json.Unmarshal(data, &logMessage); err != nil {
		fmt.Printf("Bad Message: %s\n", err)
		return
	}

	for _, l := range loggers {
		if l.Check(logMessage) {
			l.Destination.PrintMessage(logMessage)
		}
	}
}
