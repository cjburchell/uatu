package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/cjburchell/pubsub"

	uatu "github.com/cjburchell/uatu-go"
	"github.com/cjburchell/uatu/config"
	"github.com/cjburchell/uatu/loggers"
	"github.com/cjburchell/uatu/settings"
	appSettings "github.com/cjburchell/settings-go"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// IProcessor Interface
type IProcessor interface {
	Start(config settings.AppConfig)
	Stop()
}

type processor struct {
	processors   []loggers.Logger
	pubSub       pubsub.IPubSub
	subscription pubsub.ISubscription
}

func (p *processor) Stop() {
	if p.subscription != nil {
		err := p.subscription.Close()
		if err != nil {
			log.Printf("Error stopping subscription %s", err.Error())
		}
	}

}

// Start the processor
func (p *processor) Start(config settings.AppConfig) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	if config.UsePubSub {
		var err error
		p.pubSub, p.subscription, err = p.setupPubSub(config.PubSub)
		if err != nil {
			log.Printf("unable to connect to pub sub %s", err.Error())
		}
	}

	if config.UseRest {
		go func() {

			r := mux.NewRouter()
			r.Use(func(handler http.Handler) http.Handler {
				return tokenMiddleware(handler, config)
			})
			r.HandleFunc("/log", p.handelLog).Methods("POST")

			srv := &http.Server{
				Handler:      r,
				Addr:         ":" + strconv.Itoa(config.RestPort),
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			log.Printf("Handling HTTP log messages on port %d", config.RestPort)
			if err := srv.ListenAndServe(); err != nil {
				log.Printf("Error processing HTTP %s", err.Error())
			}

			wg.Done()
		}()
	}

	wg.Wait()
}

func tokenMiddleware(next http.Handler, config settings.AppConfig) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		auth := request.Header.Get("Authorization")
		if auth != "APIKEY "+config.RestToken {
			response.WriteHeader(http.StatusUnauthorized)

			log.Printf("Unauthorized %s != %s", auth, config.RestToken)
			return
		}

		next.ServeHTTP(response, request)
	})
}

func (p processor) handelLog(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var logMessage uatu.Message
	if err := decoder.Decode(&logMessage); err != nil {
		fmt.Println("Unmarshal Failed " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	go func() {
		if err := p.handleMessage(logMessage); err != nil {
			log.Printf("Unable to process log %s", err.Error())
		}
	}()
}

// Load the processors
func Load(settings appSettings.ISettings) (IProcessor, error) {
	var err error
	processors, err := loggers.Load(settings)
	if err != nil {
		return nil, err
	}

	if len(processors) == 0 {
		log.Printf("No loggers available")
		consoleLogger := loggers.Logger{Logger: config.Logger{DestinationType: "console"}}
		consoleLogger.SetMaxLevel(uatu.INFO.Severity)
		err = consoleLogger.UpdateDestination()
		if err != nil {
			return nil, err
		}

		processors = append(processors, consoleLogger)
	}

	return &processor{
		processors: processors,
	}, nil
}

func (p processor) handleMessage(logMessage uatu.Message) error {
	for _, l := range p.processors {
		if l.Check(logMessage) {
			if err := l.Destination.PrintMessage(logMessage); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p processor) setupPubSub(settings pubsub.Settings) (pubsub.IPubSub, pubsub.ISubscription, error) {

	ps, err := pubsub.Create(context.Background(), settings)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	sub, err := ps.Subscribe(context.Background(), "logs", func(data []byte) {
		logMessage := uatu.Message{}
		if err = json.Unmarshal(data, &logMessage); err != nil {
			log.Printf("Bad Message %s", err.Error())
			return
		}

		if err = p.handleMessage(logMessage); err != nil {
			log.Printf("Unable to process log %s", err.Error())
			return
		}
	})
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return ps, sub, errors.WithStack(err)
}
