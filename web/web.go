package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cjburchell/tools-go"
	"github.com/cjburchell/yasls/web/routes/logger"
	"github.com/cjburchell/yasls/web/routes/login"
	"github.com/gorilla/mux"
)

func handelStatus(w http.ResponseWriter, _ *http.Request) {
	reply, err := json.Marshal("Ok")
	if err != nil {
		fmt.Print(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(reply)
	if err != nil {
		fmt.Print(err.Error())
	}
}

// StartHTTP Service
func StartHTTP() {
	username := tools.GetEnv("ADMIN_USER", "admin")
	password := tools.GetEnv("ADMIN_PASSWORD", "admin")
	port := tools.GetEnvInt("PORTAL_PORT", 8080)

	r := mux.NewRouter()

	r.HandleFunc("/@status", handelStatus).Methods("GET")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "ui/yasls/dist/yasls/index.html")
	})

	// setup routes
	logger.SetupRoute(r)
	login.SetupRoutes(r, username, password)

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("ui/yasls/dist/yasls"))))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + strconv.Itoa(port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Hosting UI on port %d\n", port)
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}
}
