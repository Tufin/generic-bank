package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/tufin/generic-bank/auth-proxy/app"
	"github.com/tufin/generic-bank/auth-proxy/auth"
	"github.com/tufin/generic-bank/common"
)

func main() {

	// common.Auth0()
	app.Init()

	router := mux.NewRouter()
	router.HandleFunc("/callback", auth.HandleCallback)
	router.HandleFunc("/login", auth.HandleLogin)
	router.HandleFunc("/logout", auth.HandleLogout)

	proxy := auth.CreateAuthProxy()
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth.IsAuthenticate(r) {
			proxy.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
			if _, err := w.Write([]byte("Not Found :(")); err != nil {
				log.Errorf("failed to stream response with '%v'", err)
			}
		}
	})

	common.Serve(router)
}
