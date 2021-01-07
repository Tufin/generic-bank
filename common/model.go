package common

import "github.com/pborman/uuid"

type Balance struct {
	Label  string `json:"label" form:"label" binding:"required"`
	Amount int    `json:"amount" form:"amount" binding:"required"`
}

type BalanceResponse struct {
	Balance    []Balance `json:"balance"`
	CreditCard string    `json:"credit_card"`
}

type CreditCardAccount struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	LastName   string `json:"last_name"`
	CreditCard string `json:"credit_card"`
}

type Account struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
}

type RequestAccount struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
}

type AccountList struct {
	Accounts []Account `json:"accounts"`
}

func NewAccount(account RequestAccount) *Account {

	return &Account{
		ID:       uuid.New(),
		Name:     account.Name,
		LastName: account.LastName,
	}
}
