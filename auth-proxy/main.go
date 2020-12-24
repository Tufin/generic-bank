package main

import (
	"github.com/gorilla/mux"
	"github.com/tufin/generic-bank/auth-proxy/app"
	"github.com/tufin/generic-bank/auth-proxy/auth"
	"github.com/tufin/generic-bank/common"
)

func main() {

	common.Auth0()
	app.Init()

	router := mux.NewRouter()
	router.HandleFunc("/callback", auth.HandleCallback)
	router.HandleFunc("/login", auth.HandleLogin)
	router.HandleFunc("/logout", auth.HandleLogout)

	proxy := auth.CreateAuthProxy()
	router.PathPrefix("/").HandlerFunc(proxy.ServeHTTP)

	common.Serve(router)
}
