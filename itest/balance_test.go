package itest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/tufin/generic-bank/common"
)

func TestAccountBalance(t *testing.T) {

	balanceURL := CustomerAccountBalanceURL()
	log.Infof("getting account balance... '%s'", balanceURL)
	response, err := http.Get(balanceURL)
	require.NoError(t, err, balanceURL)
	defer common.CloseWithErrLog(response.Body)

	payload, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	var balance common.BalanceResponse
	require.NoError(t, json.Unmarshal(payload, &balance))
	log.Infof("Balance: '%+v'", balance)

	require.NotEmpty(t, balance.CreditCard)

	require.Len(t, balance.Balance, 3)

	require.True(t, balance.Balance[0].Amount > 0)
	require.True(t, balance.Balance[1].Amount > 0)
	require.True(t, balance.Balance[2].Amount > 0)

	require.NotEmpty(t, balance.Balance[0].Label)
	require.NotEmpty(t, balance.Balance[1].Label)
	require.NotEmpty(t, balance.Balance[2].Label)
}
