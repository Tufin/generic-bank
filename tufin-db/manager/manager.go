package manager

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/tufin/generic-bank/common"
)

type DBManager struct {
	data []common.Account
}

func NewDBManager() *DBManager {

	return &DBManager{}
}

func (dbm *DBManager) GetAccounts() []common.Account {

	log.Infof("getting DB accounts... '%+v'", dbm.data)
	return dbm.data
}

func (dbm *DBManager) Clear() {

	log.Infof("deleting DB accounts... '%+v'", dbm.data)
	dbm.data = nil
}

func (dbm *DBManager) AddAccount(account common.Account) error {

	if account.ID == "" {
		msg := fmt.Sprintf("invalid account with empty ID '%+v'", account)
		log.Infof(msg)
		return errors.New(msg)
	}
	if account.Name == "" {
		msg := fmt.Sprintf("invalid account with empty name '%+v'", account)
		log.Infof(msg)
		return errors.New(msg)
	}
	if account.ID == "" {
		msg := fmt.Sprintf("invalid account with empty last-name '%+v'", account)
		log.Infof(msg)
		return errors.New(msg)
	}

	log.Infof("adding account... '%+v'", account)
	dbm.data = append(dbm.data, account)

	return nil
}
