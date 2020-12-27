package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tufin/generic-bank/common"
)

type AccountList struct {
	Accounts []common.Account `json:"accounts"`
}

func main() {

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	sleepTime := getSleepTime()
	postgres := getPostgresAccountsUrl()
	redisURL := getRedisAccountsUrl()
	redisDeleteAccountsRequest, err := http.NewRequest(http.MethodDelete, redisURL, nil)
	if err != nil {
		log.Fatalf("failed to create http request for redis delete accounts request '%s' with '%v'", redisURL, err)
	}

	go func() {
		log.Info("Generic Bank Indexer started")
		for {
			log.Debug("Fetching and deleting Redis accounts...")
			tr := &http.Transport{}
			client := &http.Client{Transport: tr}
			response, err := client.Do(redisDeleteAccountsRequest)
			if err != nil {
				log.Errorf("failed to fetch and delete accounts from redis with '%v'", err)
			} else {
				var list AccountList
				if err := json.Unmarshal(streamToBytes(response.Body), &list); err != nil {
					log.Errorf("failed to unmarshal account list with '%v'", err)
				} else {
					insertAccounts(postgres, list)
				}
				if err := response.Body.Close(); err != nil {
					log.Errorf("failed to close redis response body with '%v'", err)
				}
			}
			time.Sleep(sleepTime)
		}
	}()

	<-stop
	log.Info("Generic Bank Indexer has been stopped")
}

func getSleepTime() time.Duration {

	sleep := getEnvVarWithDefault("SLEEP_TIME", "30s")
	ret, err := time.ParseDuration(sleep)
	if err != nil {
		log.Fatalf("invalid sleep time duration '%s' with '%+v'", sleep, err)
	}

	return ret
}

func insertAccounts(postgres string, list AccountList) error {

	for _, currAccount := range list.Accounts {
		payload, err := json.Marshal(currAccount)
		if err != nil {
			log.Errorf("failed to marshal account '%+v' with '%v'", currAccount, err)
			return err
		}
		log.Infof("inserting account to postgres... %+v", currAccount)
		response, err := http.Post(postgres, common.ContentTypeApplicationJSON, bytes.NewReader(payload))
		if err != nil {
			log.Errorf("failed to insert account '%s' to postgres (%s) with '%v'", currAccount, postgres, err)
		} else {
			if response.StatusCode != http.StatusCreated {
				log.Errorf("failed to insert account '%s' to postgres with status '%s'", currAccount, response.Status)
			} else {
				log.Infof("account '%s' created in postgres", currAccount)
			}
			if err := response.Body.Close(); err != nil {
				log.Errorf("failed to close postgres response body with '%v'", err)
			}
		}
	}

	return nil
}

func getPostgresAccountsUrl() string {

	return fmt.Sprintf("%s/accounts", getEnvVarWithDefault("POSTGRES", "http://localhost:8088"))
}

func getRedisAccountsUrl() string {

	return fmt.Sprintf("%s/accounts", getEnvVarWithDefault("REDIS", "http://localhost:8088"))
}

func getEnvVarWithDefault(varname string, defaultValue string) string {

	ret := os.Getenv(varname)
	if ret == "" {
		ret = defaultValue
	}
	log.Infof("ENV VAR '%s'='%s'", varname, ret)

	return ret
}

func streamToBytes(stream io.Reader) []byte {

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(stream); err != nil {
		log.Errorf("failed to read from stream with '%v'", err)
		return []byte{}
	}

	return buf.Bytes()
}

func getLogLevel() log.Level {

	ret := log.InfoLevel
	level := os.Getenv("LOG_LEVEL")
	if level != "" {
		if strings.EqualFold(level, "debug") {
			ret = log.DebugLevel
		} else if strings.EqualFold(level, "warn") {
			ret = log.WarnLevel
		} else if strings.EqualFold(level, "error") {
			ret = log.ErrorLevel
		}
	}

	return ret
}
