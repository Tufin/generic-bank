package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
)

var redis string

//var pgClient common.PostgresClient

func main() {

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	mode := os.Getenv("MODE")
	if mode == "admin" {
		log.Info("Admin mode")
		//pgClient = common.NewPostgresClient()
		//defer pgClient.Close()
	} else if mode == "maintenance" {
		log.Info("Maintenance mode")
	} else {
		mode = "customer"
		log.Info("Customer mode")
	}

	redis = getRedisUrl()

	router := mux.NewRouter()
	router.HandleFunc("/boa/admin/accounts", getAccounts).Methods(http.MethodGet)
	router.HandleFunc("/accounts/{account-id}", createAccount).Methods(http.MethodPost)
	router.HandleFunc("/time", getTime).Methods(http.MethodGet)

	if mode == "admin" {
		log.Info("going into admin mode")
		router.PathPrefix("/boa/admin").Handler(http.StripPrefix("/boa/admin", http.FileServer(http.Dir("/boa/html/admin/"))))
	} else if mode == "maintenance" {
		log.Info("going into maintenance mode")
		router.PathPrefix("/boa/admin").Handler(http.StripPrefix("/boa/admin", http.FileServer(http.Dir("/boa/html/maintenance/"))))
	} else if mode == "customer" {
		log.Info("going into customer mode")
		router.PathPrefix("/boa/customer").Handler(http.StripPrefix("/boa/customer", http.FileServer(http.Dir("/boa/html/customer/"))))
	}

	router.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir("/boa/html/"))))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	go func() {
		log.Info("Generic Bank Server listening on port 8085")
		if err := http.ListenAndServe(":8085", router); err != nil {
			log.Error("Generic Bank Server interrupted. ", err)
		}
	}()

	<-stop // wait for SIGINT
	log.Info("Generic Bank Server has been stopped")
}

func getRedisUrl() string {

	ret := os.Getenv("REDIS")
	if ret == "" {
		ret = "http://redis"
	}
	log.Infof("Redis: %s", ret)

	return ret
}

func getTime(w http.ResponseWriter, r *http.Request) {
	timeService := getTimeServiceUrl(r.FormValue("zone"))

	log.Infof("getting time from (%s)...", timeService)
	response, err := http.Get(timeService)
	if err != nil {
		log.Errorf("failed to get time with '%v'", err)
	} else {
		if response.StatusCode != http.StatusOK {
			log.Errorf("failed to get time with status '%s'", response.Status)
		} else {
			log.Infof("time retrieved successfully")
			w.WriteHeader(http.StatusOK)

			defer response.Body.Close()
			body, err := ioutil.ReadAll(response.Body)

			if err != nil {
				log.Errorf("failed to read time body from time service '%v'", err)
			} else {
				w.Write(body)
			}
		}
	}
}

func getTimeServiceUrl(zone string) string {

	ret := os.Getenv("TIME")
	if ret == "" {
		ret = "http://time"
	}
	ret += fmt.Sprintf("/time?zone=%s", zone)
	log.Infof("time: %s", ret)

	return ret
}

func getAccounts(w http.ResponseWriter, _ *http.Request) {

	postgres := getPostgresAccountsUrl()

	log.Infof("getting accounts from postgres (%s)...", postgres)
	response, err := http.Get(postgres)
	if err != nil {
		log.Errorf("failed to get accounts from postgres (%s) with '%v'", postgres, err)
	} else {
		if response.StatusCode != http.StatusOK {
			log.Errorf("failed to get accounts from postgres with status '%s'", response.Status)
		} else {
			log.Infof("accounts retrieved successfully")
			w.WriteHeader(http.StatusOK)

			defer response.Body.Close()
			body, err := ioutil.ReadAll(response.Body)

			if err != nil {
				log.Errorf("failed to read response body from postgres with '%v'", err)
			} else {
				w.Write(body)
			}
		}
	}
}

func getPostgresAccountsUrl() string {

	ret := os.Getenv("POSTGRES")
	if ret == "" {
		ret = "http://localhost:8088"
	}
	ret += "/accounts"
	log.Infof("Postgres: %s", ret)

	return ret
}

func createAccount(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["account-id"]
	url := fmt.Sprintf("%s/accounts", redis)
	log.Infof("Creating account '%s' in redis (%s)", id, url)
	response, err := http.Post(url, "text/plain", bytes.NewReader([]byte(id)))
	if err != nil {
		log.Errorf("Failed to add key '%s' to redis with '%v'", id, err)
		w.WriteHeader(http.StatusInternalServerError)
	} else if response.StatusCode != http.StatusCreated {
		log.Errorf("Failed to add key '%s' to redis with status '%s' (%#v)", id, response.Status, response)
		w.WriteHeader(response.StatusCode)
	} else {
		msg := fmt.Sprintf("Account '%s' has been added to redis", id)
		log.Infof(msg)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, msg)
	}
}
