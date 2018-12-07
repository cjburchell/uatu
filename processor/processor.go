package processor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/cjburchell/uatu/settings"

	"github.com/cjburchell/go-uatu"
	"github.com/cjburchell/uatu/config"
	"github.com/cjburchell/uatu/loggers"
	"github.com/gorilla/mux"
	"github.com/nats-io/go-nats"
)

var processors []loggers.Logger

// Start the processor
func Start() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	if settings.UseNats {

		log.Printf("Connecting to nats: %s", settings.NatsURL)
		var err error
		natsConn, err = setupNats(settings.NatsURL)
		if err != nil {
			log.Errorf(err, "unable to connect to nats")
		}
	}

	if settings.UseRest {
		go func() {

			r := mux.NewRouter()
			r.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {

				decoder := json.NewDecoder(r.Body)
				var logMessage log.Message
				if err := decoder.Decode(&logMessage); err != nil {
					fmt.Println("Unmarshal Failed " + err.Error())
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusCreated)

				go func() {
					if err := handleMessage(logMessage); err != nil {
						log.Error(err, "Unable to process log", err)
					}
				}()
			}).Methods("POST")

			srv := &http.Server{
				Handler:      r,
				Addr:         ":" + strconv.Itoa(settings.RestPort),
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			log.Printf("Handling HTTP log messages on port %d", settings.RestPort)
			if err := srv.ListenAndServe(); err != nil {
				log.Error(err, "Error processing HTTP")
			}

			wg.Done()
		}()
	}

	wg.Wait()
}

// Load the processors
func Load() error {
	var err error
	processors, err = loggers.Load()
	if err != nil {
		return err
	}

	if len(processors) == 0 {
		log.Warn("No loggers available")
		consoleLogger := loggers.Logger{Logger: config.Logger{DestinationType: "console"}}
		consoleLogger.SetMaxLevel(log.INFO.Severity)
		err = consoleLogger.UpdateDestination()
		if err != nil {
			return err
		}

		processors = append(processors, consoleLogger)
	}

	return nil
}

var natsConn *nats.Conn

func handleMessage(logMessage log.Message) error {
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
		logMessage := log.Message{}
		if err = json.Unmarshal(data, &logMessage); err != nil {
			log.Error(err, "Bad Message")
			return
		}

		if err = handleMessage(logMessage); err != nil {
			log.Error(err, "Unable to process log")
			return
		}
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = natsConn.Flush()
	return natsConn, errors.WithStack(err)
}
