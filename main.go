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
	err := setupNats()
	if err != nil {
		fmt.Printf("Can't connect to nats: %v\n", err)
		return
	}
	web.StartHttp()
}

var loggers = []logger.Logger{
	logger.CreateFile(),
	logger.CreateConsole(),
}

var natsConn *nats.Conn

func setupNats() error {
	natsUrl := tools.GetEnv("NATS_URL", "tcp://nats:4222")

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
		l.PrintMessage(logMessage)
	}
}
