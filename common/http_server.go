package common

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func CloudRunServeHTTP(handle func(http.ResponseWriter, *http.Request)) {

	router := mux.NewRouter()
	router.HandleFunc("/", handle).Methods(http.MethodPost)

	Serve(router)
}

func Serve(router *mux.Router) {

	ServeMulti([]*mux.Router{router}, []string{"8080"})
}

func ServeMulti(routers []*mux.Router, ports []string) {

	logVersion()

	var servers []*http.Server
	for i := 0; i < len(ports); i++ {
		servers = append(servers, &http.Server{
			Addr: fmt.Sprintf("%s:%s", "0.0.0.0", ports[i]),
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      routers[i],
		})
		go func(server *http.Server, port string) {
			log.Infof("listening on port '%s'", port)
			if err := server.ListenAndServe(); err != nil {
				log.Error(err)
			}
		}(servers[i], ports[i])
	}
	c := make(chan os.Signal, 1)
	// Graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	<-c

	for _, currServer := range servers {
		shutdown(currServer)
	}

	log.Info("exit")
	os.Exit(0)
}

// Set CORS headers for the preflight request
// CORSEnabledFunction is an example of setting CORS headers.
// For more information about CORS and CORS preflight requests, see
// https://developer.mozilla.org/en-US/docs/Glossary/Preflight_request.
func EnableCORS(w http.ResponseWriter, _ *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.WriteHeader(http.StatusNoContent)
}

func shutdown(server *http.Server) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("failed to shoutdown server with '%v'", err)
	}
}

func logVersion() {

	log.Infof("%s/%s, %s", runtime.GOOS, runtime.GOARCH, runtime.Version())
}
