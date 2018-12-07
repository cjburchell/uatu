package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cjburchell/go-uatu"

	"github.com/cjburchell/uatu/settings"

	"github.com/cjburchell/uatu/web/routes/logger"
	"github.com/cjburchell/uatu/web/routes/login"

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

	r := mux.NewRouter()

	r.HandleFunc("/@status", handelStatus).Methods("GET")
	/*r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "ui/uatu/dist/uatu/index.html")
	})*/

	// setup routes
	logger.SetupRoute(r)
	login.SetupRoutes(r, settings.PortalUsername, settings.PortalPassword)

	//r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("ui/uatu/dist/uatu"))))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + strconv.Itoa(settings.PortalPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Hosting UI on port %d", settings.PortalPort)
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}
}
