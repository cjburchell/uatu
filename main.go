package main

import (
	"encoding/json"
	"fmt"
	"github.com/cjburchell/yasls/loggers"
	"github.com/cjburchell/tools-go"
	"github.com/cjburchell/yasls-client-go"
	"github.com/cjburchell/yasls/web"
	"github.com/nats-io/go-nats"
)

func main() {
	natsUrl := tools.GetEnv("NATS_URL", "tcp://nats:4222")
	configFile := tools.GetEnv("CONFIG_FILE", "/config.json")
	//username := tools.GetEnv("ADMIN_USER", "admin")
	//password := tools.GetEnv("ADMIN_PASSWORD", "admin")

	var err error
	processors, err = loggers.Load(configFile)
	if err != nil{
		log.Print(err.Error())
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
			destConfig := json.RawMessage(data)
			consoleLogger := loggers.Logger{DestinationType: "file", DestConfig: &destConfig}
			consoleLogger.SetMaxLevel(minLevelFile)
			processors = append(processors, consoleLogger)
		}

		minLevelConsole := tools.GetEnvInt("LOG_CONSOLE_LEVEL", log.INFO.Severity)
		enabledConsole := tools.GetEnvBool("LOG_CONSOLE", false)
		if enabledConsole {
			consoleLogger := loggers.Logger{DestinationType: "console"}
			consoleLogger.SetMaxLevel(minLevelConsole)
			processors = append(processors, consoleLogger)
		}
	}

	for _, l := range processors{
		l.UpdateDestination()
	}

	natsConn, err := setupNats(natsUrl)
	if err != nil {
		fmt.Printf("Can't connect to nats: %v\n", err)
		return
	}

	defer natsConn.Close()
	web.StartHttp()
}

var processors []loggers.Logger

func setupNats(natsUrl string) (*nats.Conn, error) {
	fmt.Println("Connecting to Nats " + natsUrl)
	natsConn, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, err
	}

	_, err = natsConn.Subscribe("logs", handleMessage)
	natsConn.Flush()
	return natsConn, err
}

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
