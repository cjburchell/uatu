package processor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/cjburchell/tools-go"
	"github.com/cjburchell/yasls-client-go"
	"github.com/cjburchell/yasls/config"
	"github.com/cjburchell/yasls/loggers"
	"github.com/gorilla/mux"
	"github.com/nats-io/go-nats"
)

var processors []loggers.Logger

// Start the processor
func Start() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	useNats := tools.GetEnvBool("USE_NATS", true)
	if useNats {
		natsURL := tools.GetEnv("NATS_URL", "tcp://nats:4222")
		fmt.Printf("Connecting to nats: %s\n", natsURL)
		var err error
		natsConn, err = setupNats(natsURL)
		if err != nil {
			fmt.Printf("unable to connect to nats: %s", err)
		}
	}

	useRest := tools.GetEnvBool("USE_REST", false)
	if useRest {
		go func() {
			port := tools.GetEnvInt("REST_PORT", 8081)

			r := mux.NewRouter()
			r.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {

				decoder := json.NewDecoder(r.Body)
				var logMessage log.LogMessage
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
		maxAge := tools.GetEnvInt("LOG_FILE_MAX_AGE", 1) //days
		maxBackups := tools.GetEnvInt("LOG_FILE_MAX_BACKUPS", 20)
		maxSize := tools.GetEnvInt("LOG_FILE_MAX_SIZE", 100) * 1024 * 1024. // megabytes
		filename := tools.GetEnv("LOG_FILE_PATH", "/logs/server.log")

		minLevelFile := tools.GetEnvInt("LOG_FILE_LEVEL", log.INFO.Severity)
		enabledFile := tools.GetEnvBool("LOG_FILE", false)

		if enabledFile {
			data, _ := json.Marshal(struct {
				MaxAge     int    `json:"max_age"`
				MaxBackups int    `json:"max_backups"`
				MaxSize    int    `json:"max_size"`
				Filename   string `json:"filename"`
			}{
				MaxAge:     maxAge,
				MaxBackups: maxBackups,
				MaxSize:    maxSize,
				Filename:   filename,
			})
			destinationConfig := json.RawMessage(data)
			consoleLogger := loggers.Logger{Logger: config.Logger{DestinationType: "file", DestinationConfig: &destinationConfig}}
			consoleLogger.SetMaxLevel(minLevelFile)
			err = consoleLogger.UpdateDestination()
			if err != nil {
				return err
			}
			processors = append(processors, consoleLogger)
		}

		minLevelConsole := tools.GetEnvInt("LOG_CONSOLE_LEVEL", log.INFO.Severity)
		enabledConsole := tools.GetEnvBool("LOG_CONSOLE", false)
		if enabledConsole {
			consoleLogger := loggers.Logger{Logger: config.Logger{DestinationType: "console"}}
			consoleLogger.SetMaxLevel(minLevelConsole)
			err = consoleLogger.UpdateDestination()
			if err != nil {
				return err
			}

			processors = append(processors, consoleLogger)
		}
	}

	return nil
}

var natsConn *nats.Conn

func handleMessage(logMessage log.LogMessage) error {
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
		logMessage := log.LogMessage{}
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
