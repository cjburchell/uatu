package processor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	uatu "github.com/cjburchell/uatu-go"
	"github.com/cjburchell/uatu/config"
	"github.com/cjburchell/uatu/loggers"
	"github.com/cjburchell/uatu/settings"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

var processors []loggers.Logger

// Start the processor
func Start(config settings.AppConfig) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	if config.UseNats {

		log.Printf("Connecting to nats: %s", config.NatsURL)
		var err error
		natsConn, err = setupNats(config.NatsURL)
		if err != nil {
			log.Printf("unable to connect to nats %s", err.Error())
		}
	}

	if config.UseRest {
		go func() {

			r := mux.NewRouter()
			r.Use(func(handler http.Handler) http.Handler {
				return tokenMiddleware(handler, config)
			})
			r.HandleFunc("/log", handelLog).Methods("POST")

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

func handelLog(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var logMessage uatu.Message
	if err := decoder.Decode(&logMessage); err != nil {
		fmt.Println("Unmarshal Failed " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	go func() {
		if err := handleMessage(logMessage); err != nil {
			log.Printf("Unable to process log %s", err.Error())
		}
	}()
}

// Load the processors
func Load() error {
	var err error
	processors, err = loggers.Load()
	if err != nil {
		return err
	}

	if len(processors) == 0 {
		log.Printf("No loggers available")
		consoleLogger := loggers.Logger{Logger: config.Logger{DestinationType: "console"}}
		consoleLogger.SetMaxLevel(uatu.INFO.Severity)
		err = consoleLogger.UpdateDestination()
		if err != nil {
			return err
		}

		processors = append(processors, consoleLogger)
	}

	return nil
}

var natsConn *nats.Conn

func handleMessage(logMessage uatu.Message) error {
	for _, l := range processors {
		if l.Check(logMessage) {
			if err := l.Destination.PrintMessage(logMessage); err != nil {
				return err
			}
		}
	}

	return nil
}

func setupNats(natsURL string) (*nats.Conn, error) {
	var err error
	natsConn, err = nats.Connect(natsURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = natsConn.Subscribe("logs", func(msg *nats.Msg) {
		data := msg.Data
		logMessage := uatu.Message{}
		if err = json.Unmarshal(data, &logMessage); err != nil {
			log.Printf("Bad Message %s", err.Error())
			return
		}

		if err = handleMessage(logMessage); err != nil {
			log.Printf("Unable to process log %s", err.Error())
			return
		}
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = natsConn.Flush()
	return natsConn, errors.WithStack(err)
}
