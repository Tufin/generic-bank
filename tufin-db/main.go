package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/tufin/generic-bank/common"
	"github.com/tufin/generic-bank/tufin-db/manager"
	sabik "github.com/tufin/sabik/client"
)

var dbManager *manager.DBManager

func main() {

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	dbManager = manager.NewDBManager()

	router := mux.NewRouter()
	middleware := sabik.CreateMiddleware()
	router.Handle("/accounts", middleware.Handle(http.HandlerFunc(getAccounts))).Methods(http.MethodGet)
	router.Handle("/accounts", middleware.Handle(http.HandlerFunc(addAccount))).Methods(http.MethodPost)
	router.Handle("/accounts", middleware.Handle(http.HandlerFunc(getDeleteAccounts))).Methods(http.MethodDelete)
	router.Handle("/", middleware.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "Generic Bank DB")
		if err != nil {
			log.Errorf("failed to write response with '%v'", err)
		}
	}))).Methods(http.MethodGet)
	go func() {
		log.Info("Generic Bank DB listening on port 8088")
		if err := http.ListenAndServe(":8088", router); err != nil {
			log.Error("Generic Bank DB interrupted with ", err)
		}
	}()

	<-stop // wait for SIGINT
	log.Info("Generic Bank DB has been stopped")
}

func getAccounts(w http.ResponseWriter, r *http.Request) {

	common.RespondWith(w, r, http.StatusOK, dbManager.GetAccounts())
}

func getDeleteAccounts(w http.ResponseWriter, r *http.Request) {

	common.RespondWith(w, r, http.StatusOK, dbManager.GetAccounts())
	dbManager.Clear()
}

func addAccount(w http.ResponseWriter, r *http.Request) {

	var account common.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Infof("failed to unmarshal account with '%v'", err)
		return
	}
	common.CloseWithErrLog(r.Body)

	if err := dbManager.AddAccount(account); err != nil {
		common.RespondWith(w, r, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	common.RespondWith(w, r, http.StatusCreated, nil)
}
