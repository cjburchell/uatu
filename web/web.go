package web

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"github.com/cjburchell/yasls/web/routes"
	"github.com/cjburchell/tools-go"
)

func handelStatus(w http.ResponseWriter, _ *http.Request) {
	reply, _ := json.Marshal("Ok")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}

func StartHttp() {
	username := tools.GetEnv("ADMIN_USER", "admin")
	password := tools.GetEnv("ADMIN_PASSWORD", "admin")

	r := mux.NewRouter()

	r.HandleFunc("/@status", handelStatus).Methods("GET")
	routes.SetupLoggerRoute(r, username, password)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8088",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf(err.Error())
	}
}
