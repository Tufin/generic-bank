package common

import "github.com/pborman/uuid"

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

type RequestAccount struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
}

type Account struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
}

func NewAccount(account RequestAccount) *Account {

	return &Account{
		ID:       uuid.New(),
		Name:     account.Name,
		LastName: account.LastName,
	}
}
