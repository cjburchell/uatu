package logger

import (
	"encoding/json"
	"net/http"

	"github.com/cjburchell/uatu/config"
	"github.com/cjburchell/uatu/web/middelware"
	"github.com/gorilla/mux"
)

// SetupRoute for the logger
func SetupRoute(r *mux.Router) {
	loggerRoute := r.PathPrefix("/logger").Subrouter()
	loggerRoute.Handle("/", middelware.ValidateJWT(handleGetLoggers)).Methods("GET")
}

func handleGetLoggers(writer http.ResponseWriter, _ *http.Request) {
	loggers, _ := config.GetLoggers()
	reply, _ := json.Marshal(loggers)
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(reply)
}
