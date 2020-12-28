package manager_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/generic-bank/common"
	"github.com/tufin/generic-bank/tufin-db/manager"
)

func TestDBManager_AddAccount_EmptyID(t *testing.T) {

	require.Error(t, manager.NewDBManager().AddAccount(common.Account{
		ID:       "",
		Name:     "Israel",
		LastName: "Israeli",
	}))
}

func TestDBManager_AddAccount(t *testing.T) {

	require.NoError(t, manager.NewDBManager().AddAccount(common.Account{
		ID:       "123u",
		Name:     "Israel",
		LastName: "Israeli",
	}))
}

func TestDBManager_GetAccounts(t *testing.T) {

	const id, name, lastname = "123u", "Israel", "Israeli"

	dbManager := manager.NewDBManager()
	require.NoError(t, dbManager.AddAccount(common.Account{
		ID:       id,
		Name:     name,
		LastName: lastname,
	}))
	accounts := dbManager.GetAccounts()

	require.Len(t, accounts.Accounts, 1)
	require.Equal(t, id, accounts.Accounts[0].ID)
	require.Equal(t, name, accounts.Accounts[0].Name)
	require.Equal(t, lastname, accounts.Accounts[0].LastName)
}

func TestDBManager_ClearAccounts(t *testing.T) {

	const id, name, lastname = "1", "Israel", "Israeli"

	dbManager := manager.NewDBManager()
	require.NoError(t, dbManager.AddAccount(common.Account{
		ID:       id,
		Name:     name,
		LastName: lastname,
	}))

	require.Len(t, dbManager.GetAccounts().Accounts, 1)
	dbManager.Clear()
	require.Len(t, dbManager.GetAccounts().Accounts, 0)

	require.NoError(t, dbManager.AddAccount(common.Account{
		ID:       "2",
		Name:     name,
		LastName: lastname,
	}))
	require.Len(t, dbManager.GetAccounts().Accounts, 1)
}
