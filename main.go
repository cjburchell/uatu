package main

import (
	"log"
	"sync"

	configFile "github.com/cjburchell/settings-go"
	"github.com/cjburchell/tools-go/env"
	"github.com/cjburchell/uatu/config"
	"github.com/cjburchell/uatu/processor"
	"github.com/cjburchell/uatu/settings"
	"github.com/cjburchell/uatu/web"
)

func main() {
	wg := &sync.WaitGroup{}
	appSettings := settings.Get(configFile.Get(env.Get("SettingsFile", "")))

	log.Printf("Loading config file %s", appSettings.ConfigFile)
	err := config.Setup(appSettings.ConfigFile)
	if err != nil {
		log.Printf("Unable to load config file %s", err.Error())
	}

	log.Print("Setting up processors")
	err = processor.Load()
	if err != nil {
		log.Printf("Unable to load processors %s", err.Error())
		return
	}

	wg.Add(1)

	go func() {
		processor.Start(appSettings)
		wg.Done()
	}()

	if appSettings.PortalEnable {
		wg.Add(1)
		go func() {
			web.StartHTTP(appSettings)
			wg.Done()
		}()
	}

	wg.Wait()
}
