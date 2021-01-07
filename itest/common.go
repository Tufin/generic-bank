package itest

import (
	"fmt"

	"github.com/tufin/generic-bank/common"
)

var baseURL string

func init() {

	// os.Setenv("INGRESS_IP", "34.102.175.113")
	baseURL = fmt.Sprintf("http://%s", common.GetEnvOrExit("INGRESS_IP"))
}

func CustomerAccountURL() string {

	return fmt.Sprintf("%s/customer/accounts", baseURL)
}

func AdminAccountURL() string {

	return fmt.Sprintf("%s/admin/accounts", baseURL)
}

func CustomerAccountBalanceURL() string {

	return fmt.Sprintf("%s/customer/balance", baseURL)
}
