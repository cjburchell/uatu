package main

import (
	"log"

	"github.com/cjburchell/tools-go"
	"github.com/cjburchell/yasls/config"
	"github.com/cjburchell/yasls/processor"
	"github.com/cjburchell/yasls/web"
)

func main() {
	configFile := tools.GetEnv("CONFIG_FILE", "/config.json")

	err := config.Setup(configFile)
	if err != nil {
		log.Print(err.Error())
	}

	err = processor.Start()
	if err != nil {
		log.Print(err.Error())
	}
	defer processor.Stop()
	web.StartHttp()
}
