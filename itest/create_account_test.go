package itest

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/tufin/generic-bank/common"
)

func TestCreateAccount(t *testing.T) {

	require.True(t, isCreated(t, createAccount(t)))
}

func isCreated(t *testing.T, account common.Account) bool {

	ret := false
	for i := 0; i < 33; i++ {
		log.Infof("'%d' checking if account '%s' fetched by indexer and persisted in postgres...", i, account.ID)
		if isContain(getAllAccounts(t), account) {
			ret = true
		}
		time.Sleep(time.Second * 1)
	}

	return ret
}

func isContain(allAccounts []common.Account, account common.Account) bool {

	ret := false
	for _, currAccount := range allAccounts {
		if cmp.Equal(currAccount, account) {
			ret = true
			break
		}
	}

	return ret
}

func getAllAccounts(t *testing.T) []common.Account {

	adminAccountURL := AdminAccountURL()
	log.Infof("getting accounts... '%s'", adminAccountURL)
	response, err := http.Get(adminAccountURL)
	require.NoError(t, err, adminAccountURL)
	defer common.CloseWithErrLog(response.Body)

	payload, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	var accounts map[string][]common.Account
	require.NoError(t, json.Unmarshal(payload, &accounts))
	ret := accounts["accounts"]
	log.Infof("Accounts: %+v", ret)

	return ret
}

func createAccount(t *testing.T) common.Account {

	customerAccountURL := CustomerAccountURL()
	account := doCreateAccount()
	log.Infof("creating customer account '%+v'... '%s'", account, customerAccountURL)
	response, err := http.Post(customerAccountURL, common.ContentTypeTextPlain, getAccountRequestBody(t, account))
	require.NoError(t, err, customerAccountURL)
	defer common.CloseWithErrLog(response.Body)
	require.Equal(t, http.StatusCreated, response.StatusCode, customerAccountURL)

	payload, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	var createdAccount common.Account
	require.NoError(t, json.Unmarshal(payload, &createdAccount))
	log.Infof("account '%+v' has been created '%s'", createdAccount, customerAccountURL)
	require.Equal(t, account.Name, createdAccount.Name)
	require.Equal(t, account.LastName, createdAccount.LastName)

	return createdAccount
}

func getAccountRequestBody(t *testing.T, account common.Account) io.Reader {

	ret, err := json.Marshal(account)
	require.NoError(t, err)

	return bytes.NewReader(ret)
}

func doCreateAccount() common.Account {

	id := strconv.FormatInt(time.Now().Unix(), 10)

	return common.Account{
		Name:     "test",
		LastName: id,
	}
}
