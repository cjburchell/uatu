package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/cjburchell/tools-go"
	"github.com/nats-io/go-nats"
)

// Level of the log
type Level struct {
	// Text representation of the log
	Text string
	// Severity value of the log
	Severity int
}

var (
	// DEBUG log level
	DEBUG = Level{Text: "Debug", Severity: 0}
	// INFO log level
	INFO = Level{Text: "Info", Severity: 1}
	// WARNING log level
	WARNING = Level{Text: "Warning", Severity: 2}
	// ERROR log level
	ERROR = Level{Text: "Error", Severity: 3}
	// FATAL log level
	FATAL = Level{Text: "Fatal", Severity: 4}
)

func getStackTrace() string {

	var buffer bytes.Buffer
	_, err := buffer.WriteString(fmt.Sprintf("Stacktrace:\n"))
	if err != nil {
		return ""
	}

	i := 2
	for i < 40 {
		if function1, file1, line1, ok := runtime.Caller(i); ok {
			_, err = buffer.WriteString(fmt.Sprintf("      at %s (%s:%d)\n", runtime.FuncForPC(function1).Name(), file1, line1))
			if err != nil {
				return ""
			}
		} else {
			break
		}
		i++
	}

	return buffer.String()
}

// Warnf Print a formatted warning level message
func Warnf(format string, v ...interface{}) {
	printLog(fmt.Sprintf(format, v...), WARNING)
}

// Warn Print a warning message
func Warn(v ...interface{}) {
	printLog(fmt.Sprint(v...), WARNING)
}

// Error Print a error level message
func Error(err error, v ...interface{}) {
	msg := fmt.Sprint(v...)
	if msg == "" {
		msg = fmt.Sprintf("Error: %s\n%s", err.Error(), getStackTrace())
	} else {
		msg = fmt.Sprintf("%s\nError: %s\n%s", msg, err.Error(), getStackTrace())
	}

	printLog(msg, ERROR)
}

// Errorf Print a formatted error level message
func Errorf(err error, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if msg == "" {
		msg = fmt.Sprintf("Error: %s\n%s", err.Error(), getStackTrace())
	} else {
		msg = fmt.Sprintf("%s\nError: %s\n%s", msg, err.Error(), getStackTrace())
	}

	printLog(msg, ERROR)
}

// Fatal print fatal level message
func Fatal(v ...interface{}) {
	printLog(fmt.Sprint(v...), FATAL)
	log.Panic(v...)
}

// Fatalf print formatted fatal level message
func Fatalf(format string, v ...interface{}) {
	printLog(fmt.Sprintf(format, v...), FATAL)
	log.Panicf(format, v...)
}

// Debug print debug level message
func Debug(v ...interface{}) {
	printLog(fmt.Sprint(v...), DEBUG)
}

// Debugf print formatted debug level  message
func Debugf(format string, v ...interface{}) {
	printLog(fmt.Sprintf(format, v...), DEBUG)
}

// Print print info level message
func Print(v ...interface{}) {
	printLog(fmt.Sprint(v...), INFO)
}

// Printf print info level message
func Printf(format string, v ...interface{}) {
	printLog(fmt.Sprintf(format, v...), INFO)
}

var hostname, _ = os.Hostname()
var natsConn *nats.Conn
var restClient *http.Client

// Settings for sending logs
type Settings struct {
	ServiceName  string
	RestAddress  string
	NatsURL      string
	MinLogLevel  int
	LogToConsole bool
	UseNats      bool
	UseRest      bool
}

// CreateDefaultSettings creates a default settings object
func CreateDefaultSettings() Settings {
	var settings Settings
	settings.ServiceName = tools.GetEnv("LOG_SERVICE_NAME", "")
	settings.MinLogLevel = tools.GetEnvInt("LOG_LEVEL", INFO.Severity)
	settings.LogToConsole = tools.GetEnvBool("LOG_CONSOLE", true)
	settings.UseNats = tools.GetEnvBool("LOG_USE_NATS", true)
	settings.UseRest = tools.GetEnvBool("LOG_USE_REST", false)
	settings.RestAddress = tools.GetEnv("LOG_REST_URL", "http://logger:8082/log")
	settings.NatsURL = tools.GetEnv("LOG_NATS_URL", "tcp://nats:4222")

	return settings
}

var settings Settings

// Setup the logging system
func Setup(newSettings Settings) (err error) {
	settings = newSettings
	if settings.UseNats {
		natsConn, err = nats.Connect(settings.NatsURL)
		if err != nil {
			log.Printf("Can't connect: %v\n", err)
		}
	}

	if settings.UseRest {
		restClient = &http.Client{}
	}

	return err
}

// Message to be sent to centralized logger
type Message struct {
	Text        string `json:"text"`
	Level       Level  `json:"level"`
	ServiceName string `json:"serviceName"`
	Time        int64  `json:"time"`
	Hostname    string `json:"hostname"`
}

func (message Message) String() string {
	return fmt.Sprintf("[%s] %d %s - %s", message.Level.Text, message.Time, message.ServiceName, message.Text)
}

func printLog(text string, level Level) {
	message := Message{
		Text:        text,
		Level:       level,
		ServiceName: settings.ServiceName,
		Time:        time.Now().UnixNano() / 1000000,
		Hostname:    hostname,
	}

	if level.Severity >= settings.MinLogLevel && settings.LogToConsole {
		fmt.Println(message.String())
	}

	messageBites, err := json.Marshal(message)
	if err != nil {
		fmt.Println("error:", err)
	}
	if natsConn != nil {
		err = natsConn.Publish("logs", messageBites)
		if err != nil {
			fmt.Printf("Unable to send log to nats (%s): %s", err.Error(), message.String())
		}
	}

	if restClient != nil {
		_, err = restClient.Post(settings.RestAddress, "application/json", bytes.NewBuffer(messageBites))
		if err != nil {
			fmt.Printf("Unable to send log to %s(%s): %s", settings.RestAddress, err.Error(), message.String())
		}
	}
}
