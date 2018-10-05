package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/cjburchell/tools-go"
	"github.com/cjburchell/yasls/config"
	"github.com/cjburchell/yasls/processor"
	"github.com/cjburchell/yasls/web"
)

func main() {
	wg := &sync.WaitGroup{}

	configFile := tools.GetEnv("CONFIG_FILE", "/config.json")

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
