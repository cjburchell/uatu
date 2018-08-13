package main

import (
	"github.com/cjburchell/tools-go"
	"log"
	"github.com/cjburchell/yasls/web"
	"github.com/cjburchell/yasls/config"
	"github.com/cjburchell/yasls/processor"
)

func main() {
	configFile := tools.GetEnv("CONFIG_FILE", "/config.json")

	err := config.Setup(configFile)
	if err != nil{
		log.Print(err.Error())
	}

	err = processor.Start()
	if err != nil{
		log.Print(err.Error())
	}
	defer processor.Stop()
	web.StartHttp()
}


