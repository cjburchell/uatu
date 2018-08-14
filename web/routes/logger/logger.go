package logger

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"log"
	"github.com/cjburchell/yasls/config"
	"github.com/cjburchell/yasls/processor"
	"github.com/cjburchell/yasls/web/middelware"
)

func SetupRoute(r *mux.Router) {
	loggerRoute := r.PathPrefix("/logger").Subrouter()
	loggerRoute.Handle("/", middelware.ValidateJWT(handleGetLoggers)).Methods("GET")
	loggerRoute.Handle("/{Id}", middelware.ValidateJWT(handleGetLogger)).Methods("GET")
	loggerRoute.Handle("/{Id}", middelware.ValidateJWT(handleUpdateLogger)).Methods("POST")
	loggerRoute.Handle("/", middelware.ValidateJWT(handleAddLogger)).Methods("PUT")
	loggerRoute.Handle("/{Id}", middelware.ValidateJWT(handleDeleteLogger)).Methods("DELETE")
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