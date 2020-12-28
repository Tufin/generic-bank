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

	jwtMiddleware := auth.CreateJWTMiddleware()
	proxy := auth.CreateAuthProxy()
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !auth.IsAuthenticate(r) {
			w.WriteHeader(http.StatusNotFound)
			if _, err := w.Write([]byte("Not Found :(")); err != nil {
				log.Errorf("failed to stream response with '%v'", err)
			}
			return
		}

		if err := jwtMiddleware.CheckJWT(w, r); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte("Unauthorized :(")); err != nil {
				log.Errorf("failed to stream response with '%v'", err)
			}
			return
		}

		proxy.ServeHTTP(w, r)
	})

	common.Serve(router)
}
