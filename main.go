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
	config := configFile.Get(env.Get("SettingsFile", "config.json"))
	appSettings := settings.Get(config)

	log.Print("Setting up processors")
	p, err := processor.Load(config)
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
