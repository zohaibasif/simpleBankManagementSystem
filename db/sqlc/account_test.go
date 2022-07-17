package db

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zohaibAsif/simple_bank_management_system/util"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func createRandomAccount(t *testing.T) Account {

	user := createRandomUser(t)

	args := CreateAccountParams{
		Owner:    user.Username,
		Currency: util.RandomCurrency(),
		Balance:  util.RandomAmount(),
	}

	account, err := queries.CreateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Currency, account.Currency)
	require.Equal(t, args.Balance, account.Balance)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {

	createRandomAccount(t)

}

func TestGetAccount(t *testing.T) {

	account := createRandomAccount(t)

	response, err := queries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)

	require.NotEmpty(t, response)

	require.Equal(t, account.ID, response.ID)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	args := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := queries.ListAccounts(context.Background(), args)

	require.NoError(t, err)

	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
