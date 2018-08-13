package processor

import (
	"github.com/nats-io/go-nats"
	"github.com/cjburchell/yasls-client-go"
	"encoding/json"
	"fmt"
	"github.com/cjburchell/yasls/loggers"
	"github.com/cjburchell/tools-go"
	"github.com/cjburchell/yasls/config"
)

var processors []loggers.Logger

func Start() error {
	err := LoadProcessors()
	if err != nil{
		return err
	}

	natsUrl := tools.GetEnv("NATS_URL", "tcp://nats:4222")
	natsConn, err = setupNats(natsUrl)
	return err
}

func LoadProcessors() error  {
	var err error
	processors, err = loggers.Load()
	if err != nil{
		return err
	}

	if len(processors) == 0 {
		maxAge := tools.GetEnvInt("LOG_FILE_MAX_AGE", 1) //days
		maxBackups := tools.GetEnvInt("LOG_FILE_MAX_BACKUPS", 20)
		maxSize := tools.GetEnvInt("LOG_FILE_MAX_SIZE", 100) * 1024 * 1024. // megabytes
		filename := tools.GetEnv("LOG_FILE_PATH", "/logs/server.log")

		minLevelFile := tools.GetEnvInt("LOG_FILE_LEVEL", log.INFO.Severity)
		enabledFile := tools.GetEnvBool("LOG_FILE", false)

		if enabledFile {
			data, _ := json.Marshal(struct{
				MaxAge     int                `json:"max_age"`
				MaxBackups int                `json:"max_backups"`
				MaxSize    int                `json:"max_size"`
				Filename   string             `json:"filename"`
			}{
				MaxAge:maxAge,
				MaxBackups:maxBackups,
				MaxSize:maxSize,
				Filename:filename,
			})
			destinationConfig := json.RawMessage(data)
			consoleLogger := loggers.Logger{Logger: config.Logger{ DestinationType: "file", DestinationConfig: &destinationConfig}}
			consoleLogger.SetMaxLevel(minLevelFile)
			consoleLogger.UpdateDestination()
			processors = append(processors, consoleLogger)
		}

		minLevelConsole := tools.GetEnvInt("LOG_CONSOLE_LEVEL", log.INFO.Severity)
		enabledConsole := tools.GetEnvBool("LOG_CONSOLE", false)
		if enabledConsole {
			consoleLogger := loggers.Logger{Logger: config.Logger{DestinationType: "console"}}
			consoleLogger.SetMaxLevel(minLevelConsole)
			consoleLogger.UpdateDestination()
			processors = append(processors, consoleLogger)
		}
	}

	return nil
}

func Stop(){
	natsConn.Close()
}

var natsConn *nats.Conn

func handleMessage(message *nats.Msg) {
	data := message.Data
	logMessage := log.LogMessage{}
	if err := json.Unmarshal(data, &logMessage); err != nil {
		fmt.Printf("Bad Message: %s\n", err)
		return
	}

	for _, l := range processors {
		if l.Check(logMessage) {
			l.Destination.PrintMessage(logMessage)
		}
	}
}

func setupNats(natsUrl string) (*nats.Conn, error) {
	natsConn, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, err
	}

	_, err = natsConn.Subscribe("logs", handleMessage)
	natsConn.Flush()
	return natsConn, err
}
