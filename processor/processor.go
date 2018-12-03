package processor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/cjburchell/tools-go/env"

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

	useNats := env.GetBool("USE_NATS", true)
	if useNats {
		natsURL := env.Get("NATS_URL", "tcp://nats:4222")
		fmt.Printf("Connecting to nats: %s\n", natsURL)
		var err error
		natsConn, err = setupNats(natsURL)
		if err != nil {
			fmt.Printf("unable to connect to nats: %s", err)
		}
	}

	useRest := env.GetBool("USE_REST", false)
	if useRest {
		go func() {
			port := env.GetInt("REST_PORT", 8081)

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
						fmt.Printf("Unable to process log: %s\n", err)
					}
				}()
			}).Methods("POST")

			srv := &http.Server{
				Handler:      r,
				Addr:         ":" + strconv.Itoa(port),
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			fmt.Printf("Handling HTTP log meessges on port %d\n", port)
			if err := srv.ListenAndServe(); err != nil {
				fmt.Print(err.Error())
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
		return nil, err
	}

	_, err = natsConn.Subscribe("logs", func(msg *nats.Msg) {
		data := msg.Data
		logMessage := log.Message{}
		if err = json.Unmarshal(data, &logMessage); err != nil {
			fmt.Printf("Bad Message: %s\n", err)
			return
		}

		if err = handleMessage(logMessage); err != nil {
			fmt.Printf("Unable to process log: %s\n", err)
			return
		}
	})
	if err != nil {
		return nil, err
	}

	err = natsConn.Flush()
	return natsConn, err
}
