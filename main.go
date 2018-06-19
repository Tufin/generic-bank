package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	redis      string
	balanceURL string
)

type Balance struct {
	Label  string `json:"label" form:"label" binding:"required"`
	Amount int    `json:"amount" form:"amount" binding:"required"`
}

func main() {

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	redis = getRedisUrl()
	balanceURL = getBalanceURL()

	router := mux.NewRouter()
	router.HandleFunc("/boa/admin/accounts", getAccounts).Methods(http.MethodGet)
	router.HandleFunc("/accounts/{account-id}", createAccount).Methods(http.MethodPost)
	router.HandleFunc("/time", getTime).Methods(http.MethodGet)

	mode := os.Getenv("MODE")
	if mode == "admin" {
		log.Info("Admin mode")
		router.PathPrefix("/admin").Handler(angularRouteHandler("/admin", getAngularAssets("/boa/html/")))
		router.Handle("/admin/", angularRouteHandler("/admin", http.HandlerFunc(getAngularApp)))
	} else if mode == "balance" {
		log.Info("Balance Mode")
		router.HandleFunc("/balance", getRandomBalance).Methods(http.MethodGet)
	} else {
		log.Info("Customer Mode")
		router.HandleFunc("/balance", getBalanceAsCustomer).Methods(http.MethodGet)
		router.PathPrefix("/customer").Handler(angularRouteHandler("/customer", getAngularAssets("/boa/html/")))
		router.Handle("/customer/", angularRouteHandler("/customer", http.HandlerFunc(getAngularApp)))
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	go func() {
		const port = ":8085"
		log.Infof("Generic Bank Server listening on port %s", port)
		if err := http.ListenAndServe(port, router); err != nil {
			log.Error("Generic Bank Server interrupted. ", err)
		}
	}()

	<-stop // wait for SIGINT
	log.Info("Generic Bank Server has been stopped")
}

func getAngularAssets(path string) http.Handler {

	return http.FileServer(http.Dir(path))
}

func getAngularApp(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "/boa/html/index.html")
}

func angularRouteHandler(path string, h http.Handler) http.Handler {

	return http.StripPrefix(path, h)
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

func getBalanceAsCustomer(w http.ResponseWriter, _ *http.Request) {

	if response, err := http.Get(balanceURL); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("failed to to get balance from service with '%v'", err)
	} else {
		defer response.Body.Close()
		if body, err := ioutil.ReadAll(response.Body); err != nil {
			log.Errorf("failed to read balance body from balance service with '%v'", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Infof("Customer balance: '%s'", string(body))
			w.Write(body)
		}
	}
}

func getBalanceURL() string {

	ret := os.Getenv("BALANCE_URL")
	if ret == "" {
		ret = "http://localhost:8085"
	}
	ret += "/balance"
	log.Infof("Balance: '%s'", ret)

	return ret
}

func getRandomBalance(w http.ResponseWriter, _ *http.Request) {

	ret := []Balance{
		{
			Label:  "Savings",
			Amount: random(0, 150000),
		},
		{
			Label:  "Investments",
			Amount: random(0, 100000),
		},
		{
			Label:  "Balance",
			Amount: random(-15000, 150000),
		},
	}
	log.Infof("Random balance: '%+v'", ret)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ret)
}

func random(min, max int) int {

	rand.Seed(time.Now().Unix())

	return rand.Intn(max-min) + min
}
