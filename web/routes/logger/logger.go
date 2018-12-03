package logger

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cjburchell/uatu/config"
	"github.com/cjburchell/uatu/processor"
	"github.com/cjburchell/uatu/web/middelware"
	"github.com/gorilla/mux"
)

// SetupRoute for the logger
func SetupRoute(r *mux.Router) {
	loggerRoute := r.PathPrefix("/logger").Subrouter()
	loggerRoute.Handle("/", middelware.ValidateJWT(handleGetLoggers)).Methods("GET")
	loggerRoute.Handle("/{Id}", middelware.ValidateJWT(handleGetLogger)).Methods("GET")
	loggerRoute.Handle("/{Id}", middelware.ValidateJWT(handleUpdateLogger)).Methods("POST")
	loggerRoute.Handle("/", middelware.ValidateJWT(handleAddLogger)).Methods("PUT")
	loggerRoute.Handle("/{Id}", middelware.ValidateJWT(handleDeleteLogger)).Methods("DELETE")
}

func handleGetLoggers(writer http.ResponseWriter, _ *http.Request) {
	loggers, _ := config.GetLoggers()
	reply, _ := json.Marshal(loggers)
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(reply)
}

func handleGetLogger(writer http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loggerID := vars["Id"]
	loggers, _ := config.GetLogger(loggerID)
	reply, _ := json.Marshal(loggers)
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(reply)
}

func handleUpdateLogger(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var logger config.Logger
	err := decoder.Decode(&logger)
	if err != nil {
		log.Print("Unmarshal Failed " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	config.UpdateLogger(logger)
	processor.Load()
	w.WriteHeader(http.StatusAccepted)
}

func handleAddLogger(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var logger config.Logger
	err := decoder.Decode(&logger)
	if err != nil {
		log.Print("Unmarshal Failed " + err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	config.AddLogger(logger)
	processor.Load()
	writer.WriteHeader(http.StatusAccepted)
}

func handleDeleteLogger(writer http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loggerID := vars["Id"]

	config.DeleteLogger(loggerID)
	processor.Load()

	writer.WriteHeader(http.StatusAccepted)
}
