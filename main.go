package main

import (
	"log"
	"sync"

	configFile "github.com/cjburchell/settings-go"
	"github.com/cjburchell/tools-go/env"
	"github.com/cjburchell/uatu/processor"
	"github.com/cjburchell/uatu/settings"
)

func main() {
	wg := &sync.WaitGroup{}
	appSettings := settings.Get(configFile.Get(env.Get("SettingsFile", "settings.yml")))

	log.Printf("Loading config file %s", appSettings.ConfigFile)
	log.Print("Setting up processors")
	p, err := processor.Load(appSettings.ConfigFile)
	if err != nil {
		log.Printf("Unable to load processors %s", err.Error())
		return
	}
	defer p.Stop()

	wg.Add(1)
	go func() {
		p.Start(appSettings)
		wg.Done()
	}()

	wg.Wait()
}
