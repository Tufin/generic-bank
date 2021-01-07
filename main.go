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
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/tufin/generic-bank/common"
	sabik "github.com/tufin/sabik/client"
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
	if err := os.Setenv(sabik.EnvKeyServiceName, mode); err != nil {
		log.Fatalf("faile to set environment variable with '%v'", err)
	}
	middleware := sabik.CreateMiddleware()

	if mode == "admin" {
		router.Handle("/admin/accounts", middleware.Handle(http.HandlerFunc(handleGetAccounts))).Methods(http.MethodGet)
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

func handleGetAccounts(w http.ResponseWriter, r *http.Request) {

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
	var dbAccounts common.AccountList
	err = json.Unmarshal(body, &dbAccounts)
	if err != nil {
		log.Errorf("failed to unmarshal response body from postgres into accounts with '%v'", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	accounts := getAccounts(dbAccounts, r.FormValue("mode"))
	log.Info(accounts)

	ret, err := json.Marshal(accounts)
	if err != nil {
		log.Errorf("failed to marshal accounts with '%v'", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(ret)
	if err != nil {
		log.Errorf("failed to write accounts response with '%v'", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getAccounts(dbAccounts common.AccountList, mode string) map[string][]map[string]string {

	if mode == "compact" {
		return getAccountsWithoutCreditCards(dbAccounts)
	}

	return getAccountsWithCreditCards(dbAccounts)
}

func getAccountsWithoutCreditCards(accounts common.AccountList) map[string][]map[string]string {

	var ret []map[string]string
	for _, currAccount := range accounts.Accounts {
		ret = append(ret, accountToMap(currAccount))
	}

	return map[string][]map[string]string{"accounts": ret}
}

func accountToMap(account common.Account) map[string]string {

	return map[string]string{
		"id":        account.ID,
		"name":      account.Name,
		"last_name": account.LastName,
	}
}

var creditCards = []string{"4485281688960105", "4532343129620269", "4485 2816 8896 0105",
	"4485-2816-8896-0105", "6011128940161477", "6011241861038622", "30014800829892",
	"30507075123768", "30450567425005", "2720993263757610", "2221006966003549",
	"2720995880263559", "378373260929703",
}

func getAccountsWithCreditCards(dbAccounts common.AccountList) map[string][]map[string]string {

	count := len(creditCards)
	next := 0
	var ret []map[string]string
	for i := 0; i < len(dbAccounts.Accounts); i++ {
		ret = append(ret, accountWithCreditCardToMap(common.CreditCardAccount{
			Name:       dbAccounts.Accounts[i].Name,
			LastName:   dbAccounts.Accounts[i].LastName,
			ID:         dbAccounts.Accounts[i].ID,
			CreditCard: creditCards[next],
		}))
		next++
		if next == count {
			next = 0
		}
	}

	return map[string][]map[string]string{"accounts": ret}
}

func accountWithCreditCardToMap(account common.CreditCardAccount) map[string]string {

	return map[string]string{
		"id":          account.ID,
		"name":        account.Name,
		"last_name":   account.LastName,
		"credit_card": account.CreditCard,
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

	var requestAccount common.RequestAccount
	if err := json.NewDecoder(r.Body).Decode(&requestAccount); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Infof("failed to unmarshal account request with '%v'", err)
		return
	}
	defer common.CloseWithErrLog(r.Body)

	account := common.NewAccount(requestAccount)
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
	log.Infof("account '%+v' added to redis", account)

	common.RespondWith(w, r, http.StatusCreated, account)
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
			if _, err := w.Write(body); err != nil {
				log.Errorf("failed to write response with '%v'", err)
			}
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

	ret := common.BalanceResponse{Balance: []common.Balance{
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
	}, CreditCard: "4485281688960105"}
	log.Infof("Random balance: '%+v'", ret)
	w.Header().Set(common.HeaderContentType, common.ContentTypeApplicationJSON)
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		log.Errorf("failed to write response with '%v'", err)
	}
}

func random(min, max int) int {

	rand.Seed(time.Now().Unix())

	return rand.Intn(max-min) + min
}
