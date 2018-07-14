package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cjburchell/tools-go"
	"github.com/nats-io/go-nats"
	"log"
	"os"
	"runtime"
	"time"
)

type Level struct {
	Text     string
	Severity int
}

var (
	DEBUG   = Level{Text: "Debug", Severity: 0}
	INFO    = Level{Text: "Info", Severity: 1}
	WARNING = Level{Text: "Warning", Severity: 2}
	ERROR   = Level{Text: "Error", Severity: 3}
	FATAL   = Level{Text: "Fatal", Severity: 4}
)

func getStackTrace() string {

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Stacktrace:\n"))
	i := 2
	for i < 40 {
		if function1, file1, line1, ok := runtime.Caller(i); ok {
			buffer.WriteString(fmt.Sprintf("      at %s (%s:%d)\n", runtime.FuncForPC(function1).Name(), file1, line1))
		} else {
			break
		}
		i++
	}

	return buffer.String()
}

func Warnf(format string, v ...interface{}) {
	printLog(fmt.Sprintf(format, v...), WARNING)
}

func Warn(v ...interface{}) {
	printLog(fmt.Sprint(v...), WARNING)
}

func Error(err error, v ...interface{}) {
	msg := fmt.Sprint(v...)
	if msg == "" {
		msg = fmt.Sprintf("Error: %s\n%s", err.Error(), getStackTrace())
	} else {
		msg = fmt.Sprintf("%s\nError: %s\n%s", msg, err.Error(), getStackTrace())
	}

	printLog(msg, ERROR)
}

func Errorf(err error, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if msg == "" {
		msg = fmt.Sprintf("Error: %s\n%s", err.Error(), getStackTrace())
	} else {
		msg = fmt.Sprintf("%s\nError: %s\n%s", msg, err.Error(), getStackTrace())
	}

	printLog(msg, ERROR)
}

func Fatal(v ...interface{}) {
	printLog(fmt.Sprint(v...), FATAL)
	log.Panic(v...)
}

func Fatalf(format string, v ...interface{}) {
	printLog(fmt.Sprintf(format, v...), FATAL)
	log.Panicf(format, v...)
}

func Debug(v ...interface{}) {
	printLog(fmt.Sprint(v...), DEBUG)
}

func Debugf(format string, v ...interface{}) {
	printLog(fmt.Sprintf(format, v...), DEBUG)
}

func Print(v ...interface{}) {
	printLog(fmt.Sprint(v...), INFO)
}

func Printf(format string, v ...interface{}) {
	printLog(fmt.Sprintf(format, v...), INFO)
}

var hostname, _ = os.Hostname()
var serviceName = ""
var minLogLevel = INFO.Severity
var logToConsole = true
var natsConn *nats.Conn

func Setup() (err error) {
	serviceName = tools.GetEnv("LOG_SERVICE_NAME", "")
	minLogLevel = tools.GetEnvInt("LOG_LEVEL", INFO.Severity)
	logToConsole = tools.GetEnvBool("LOG_CONSOLE", true)
	natsUrl := tools.GetEnv("NATS_URL", "tcp://nats:4222")

	natsConn, err = nats.Connect(natsUrl)
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	return err
}

type LogMessage struct {
	Text        string
	Level       Level
	ServiceName string
	Time        int64
	Hostname    string
}

func (message LogMessage) String() string {
	return fmt.Sprintf("[%s] %d %s - %s", message.Level.Text, message.Time, message.ServiceName, message.Text)
}

func printLog(text string, level Level) {
	message := LogMessage{
		Text:        text,
		Level:       level,
		ServiceName: serviceName,
		Time:        time.Now().UnixNano() / 1000000,
		Hostname:    hostname,
	}

	if level.Severity >= minLogLevel && logToConsole {
		fmt.Println(message.String())
	}

	if natsConn != nil {
		messageBites, err := json.Marshal(message)
		if err != nil {
			fmt.Println("error:", err)
		}

		natsConn.Publish("logs", messageBites)
	}
}
