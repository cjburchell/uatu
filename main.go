package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/cjburchell/tools-go/env"

	"github.com/cjburchell/uatu/config"
	"github.com/cjburchell/uatu/processor"
	"github.com/cjburchell/uatu/web"
)

func main() {
	wg := &sync.WaitGroup{}

	configFile := env.Get("CONFIG_FILE", "/config.json")

	err := config.Setup(configFile)
	if err != nil {
		log.Print(err.Error())
	}

	err = processor.Load()
	if err != nil {
		fmt.Printf("unable to load processors: %s", err)
		return
	}

	wg.Add(2)

	go func() {
		processor.Start()
		wg.Done()
	}()
	go func() {
		web.StartHTTP()
		wg.Done()
	}()

	wg.Wait()
}
