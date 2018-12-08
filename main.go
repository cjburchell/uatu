package main

import (
	"sync"

	"github.com/cjburchell/uatu/web"

	"github.com/cjburchell/go-uatu"

	"github.com/cjburchell/uatu/settings"

	"github.com/cjburchell/uatu/config"
	"github.com/cjburchell/uatu/processor"
)

func main() {
	wg := &sync.WaitGroup{}
	err := log.Setup(log.Settings{
		ServiceName:  "logger",
		UseRest:      false,
		UseNats:      false,
		LogToConsole: true,
		MinLogLevel:  log.DEBUG,
	})

	log.Printf("Loading config file %s", settings.ConfigFile)
	err = config.Setup(settings.ConfigFile)
	if err != nil {
		log.Error(err, "Unable to load config file")
	}

	log.Print("Setting up processors")
	err = processor.Load()
	if err != nil {
		log.Error(err, "Unable to load processors")
		return
	}

	wg.Add(1)

	go func() {
		processor.Start()
		wg.Done()
	}()

	if settings.PortalEnable {
		wg.Add(1)
		go func() {
			web.StartHTTP()
			wg.Done()
		}()
	}

	wg.Wait()
}
