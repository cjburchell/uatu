package routes

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"log"
	"github.com/cjburchell/yasls/config"
	"crypto/subtle"
	"github.com/cjburchell/yasls/processor"
)

func SetupLoggerRoute(r *mux.Router, username string, password string) {
	loggerRoute := r.PathPrefix("/logger").Subrouter()
	loggerRoute.HandleFunc("/", BasicAuth(handleGetLoggers, username, password)).Methods("GET")
	loggerRoute.HandleFunc("/{Id}", BasicAuth(handleGetLogger,username, password)).Methods("GET")
	loggerRoute.HandleFunc("/{Id}", BasicAuth(handleUpdateLogger, username, password)).Methods("POST")
	loggerRoute.HandleFunc("/", BasicAuth(handleAddLogger, username, password)).Methods("PUT")
	loggerRoute.HandleFunc("/{Id}", BasicAuth(handleDeleteLogger, username, password)).Methods("DELETE")
}

func BasicAuth(handler http.HandlerFunc, username string, password string) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password for this site"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}

}


func handleGetLoggers(writer http.ResponseWriter, _ *http.Request)  {
	loggers, _ := config.GetLoggers()
	reply, _ := json.Marshal(loggers)
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(reply)
}

func handleGetLogger(writer http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	loggerId := vars["Id"]
	loggers, _ := config.GetLogger(loggerId)
	reply, _ := json.Marshal(loggers)
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(reply)
}

func handleUpdateLogger(w http.ResponseWriter, r *http.Request)  {
	decoder := json.NewDecoder(r.Body)
	var logger config.Logger
	err := decoder.Decode(&logger)
	if err != nil {
		log.Print("Unmarshal Failed " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	config.UpdateLogger(logger)
	processor.LoadProcessors()
	w.WriteHeader(http.StatusAccepted)
}

func handleAddLogger(writer http.ResponseWriter, request *http.Request)  {
	decoder := json.NewDecoder(request.Body)
	var logger config.Logger
	err := decoder.Decode(&logger)
	if err != nil {
		log.Print("Unmarshal Failed " + err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	config.AddLogger(logger)
	processor.LoadProcessors()
	writer.WriteHeader(http.StatusAccepted)
}

func handleDeleteLogger(writer http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	loggerId := vars["Id"]

	config.DeleteLogger(loggerId)
	processor.LoadProcessors()
	writer.WriteHeader(http.StatusAccepted)
}