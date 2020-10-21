package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	sabik "github.com/tufin/sabik/client"
	"github.com/tufin/sabik/common/env"
)

var (
	redis      string
	balanceURL string
)

type Balance struct {
	Label  string `json:"label" form:"label" binding:"required"`
	Amount int    `json:"amount" form:"amount" binding:"required"`
}

type SSNAccount struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	SSN      string `json:"ssn"`
}

func main() {

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	redis = getRedisUrl()
	balanceURL = getBalanceURL()
	mode := os.Getenv("MODE")

	router := mux.NewRouter()
	middleware := createMiddleware(mode)
	router.Handle("/boa/admin/accounts", middleware.Handle(http.HandlerFunc(getAccounts))).Methods(http.MethodGet)
	router.Handle("/accounts/{account-id}", middleware.Handle(http.HandlerFunc(createAccount))).Methods(http.MethodPost)
	router.Handle("/time", middleware.Handle(http.HandlerFunc(getTime))).Methods(http.MethodGet)

	if mode == "admin" {
		log.Info("Admin mode")
		router.PathPrefix("/admin/").Handler(angularRouteHandler("/admin", getAngularAssets("/boa/html/")))
	} else if mode == "balance" {
		log.Info("Balance Mode")
		router.Handle("/balance", middleware.Handle(http.HandlerFunc(getRandomBalance))).Methods(http.MethodGet)
	} else if mode == "customer" {
		log.Info("Customer Mode")
		router.Handle("/balance", middleware.Handle(http.HandlerFunc(getBalanceAsCustomer))).Methods(http.MethodGet)
		router.PathPrefix("/customer/").Handler(angularRouteHandler("/customer", getAngularAssets("/boa/html/")))
	} else {
		log.Fatalf("invalid mode '%s'", mode)
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
	log.Infof("Getting time from '%s'...", timeService)

	// DO NOT CHANGE into below http get, it's breaking the integration tests, probably because gRPC cache
	// http.Get(timeService)
	tr := &http.Transport{}
	defer tr.CloseIdleConnections()
	client := &http.Client{Transport: tr}
	getTimeRequest, err := http.NewRequest(http.MethodGet, timeService, nil)
	if err != nil {
		log.Error("failed to create get time request failed with '%v'", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := client.Do(getTimeRequest)
	// END: DO NOT CHANGE

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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response.StatusCode != http.StatusOK {
		log.Errorf("failed to get accounts from postgres with status '%s'", response.Status)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("failed to read response body from postgres with '%v'", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var dbAccounts struct {
		Accounts []string `json:"accounts"`
	}
	err = json.Unmarshal(body, &dbAccounts)
	if err != nil {
		log.Errorf("failed to unmarshal response body from postgres into accounts with '%v'", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	accounts := getAccountWithSSN(dbAccounts.Accounts)
	log.Info(accounts)

	ret, err := json.Marshal(accounts)
	if err != nil {
		log.Errorf("failed to marshal SSN accounts with '%v'", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(ret)
	if err != nil {
		log.Errorf("failed to write SSN accounts response with '%v'", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

var ssnNumbers = []string{"206-04-5678", "248-95-3456", "349-02-1234", "130-96-4321",
	"007-10-5678", "150-20-4321", "148-17-1234", "163-85-1234", "163-84-1234", "735-11-5678"}

func getAccountWithSSN(dbAccounts []string) map[string][]SSNAccount {

	count := len(ssnNumbers)
	next := 0
	var ret []SSNAccount
	for i := 0; i < len(dbAccounts); i++ {
		arr := strings.Split(dbAccounts[i], ":")
		ret = append(ret, SSNAccount{
			Name:     arr[0],
			Lastname: arr[1],
			ID:       arr[2],
			SSN:      ssnNumbers[next],
		})
		next++
		if next == count {
			next = 0
		}
	}

	return map[string][]SSNAccount{"accounts": ret}
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
			Label:  "Checking",
			Amount: random(-15000, 150000),
		},
	}
	log.Infof("Random balance: '%+v'", ret)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"balance": ret,
		"cards":   []string{"4485281688960105", "4532343129620269", "4716106401131630485"},
	})
}

func random(min, max int) int {

	rand.Seed(time.Now().Unix())

	return rand.Intn(max-min) + min
}

func createMiddleware(serviceName string) *sabik.Middleware {

	fatalOnError(os.Setenv(env.KeyDomain, "generic-bank"))
	fatalOnError(os.Setenv(env.KeyProject, "retail"))
	fatalOnError(os.Setenv(sabik.EnvKeyServiceName, serviceName))

	return sabik.NewMiddleware()
}

func fatalOnError(err error) {

	if err != nil {
		log.Fatal(err)
	}
}
