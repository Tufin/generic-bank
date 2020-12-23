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
	"github.com/tufin/generic-bank/common"
)

var (
	redis      string
	balanceURL string
)

func main() {

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	redis = getRedisUrl()
	balanceURL = getBalanceURL()
	mode := common.GetEnvOrExit("MODE")

	router := mux.NewRouter()
	middleware := common.CreateMiddleware(mode)

	if mode == "admin" {
		router.Handle("/admin/accounts", middleware.Handle(http.HandlerFunc(getAccounts))).Methods(http.MethodGet)
		router.Handle("/admin/time", middleware.Handle(http.HandlerFunc(getTime))).Methods(http.MethodGet)
		router.PathPrefix("/admin/").Handler(angularRouteHandler("/admin", getAngularAssets("/boa/html/")))
	} else if mode == "balance" {
		router.Handle("/balance", middleware.Handle(http.HandlerFunc(getRandomBalance))).Methods(http.MethodGet)
	} else if mode == "customer" {
		router.Handle("/customer/accounts", middleware.Handle(http.HandlerFunc(createAccount))).Methods(http.MethodPost)
		router.Handle("/customer/balance", middleware.Handle(http.HandlerFunc(getBalanceAsCustomer))).Methods(http.MethodGet)
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
			log.Errorf("Generic-Bank Server interrupted with '%v'", err)
		}
	}()

	<-stop // wait for SIGINT
	log.Info("Generic-Bank Server has been stopped")
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response.StatusCode != http.StatusOK {
		log.Errorf("failed to get time with status '%s'", response.Status)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer common.CloseWithErrLog(response.Body)
	log.Infof("time retrieved successfully")

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("failed to read time body from time service '%v'", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(body); err != nil {
		log.Errorf("failed to time response '%s' to stream with '%v'", string(body), err)
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
	defer common.CloseWithErrLog(response.Body)

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

func getAccountWithSSN(dbAccounts []string) map[string][]common.SSNAccount {

	count := len(ssnNumbers)
	next := 0
	var ret []common.SSNAccount
	for i := 0; i < len(dbAccounts); i++ {
		arr := strings.Split(dbAccounts[i], ":")
		ret = append(ret, common.SSNAccount{
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

	return map[string][]common.SSNAccount{"accounts": ret}
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

	var account common.RequestAccount
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Infof("failed to unmarshal account request with '%v'", err)
		return
	}
	defer common.CloseWithErrLog(r.Body)

	url := fmt.Sprintf("%s/accounts", redis)
	log.Infof("adding account '%+v' to redis '%s'", account, url)
	payload, err := json.Marshal(account)
	if err != nil {
		log.Errorf("failed to marshal account '%+v' with '%v'", account, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := http.Post(url, common.ContentTypeApplicationJSON, bytes.NewReader(payload))
	if err != nil {
		log.Errorf("failed to add account '%+v' to redis with '%v'", account, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response.StatusCode != http.StatusCreated {
		log.Errorf("failed to add account '%+v' to redis with status '%s' (%#v)", account, response.Status, response)
		w.WriteHeader(response.StatusCode)
		return
	}

	msg := fmt.Sprintf("account '%+v' added to redis", account)
	log.Infof(msg)
	w.WriteHeader(http.StatusCreated)
	if _, err := fmt.Fprint(w, msg); err != nil {
		log.Errorf("failed to stream response message '%s' with '%v'", msg, err)
	}
}

func getBalanceAsCustomer(w http.ResponseWriter, _ *http.Request) {

	if response, err := http.Get(balanceURL); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("failed to to get balance from service with '%v'", err)
	} else {
		defer common.CloseWithErrLog(response.Body)
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

	ret := []common.Balance{
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
