package db

import (
	"context"
	"database/sql"
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

func TestDeleteAccount(t *testing.T) {

	account := createRandomAccount(t)

	err := queries.DeleteAccount(context.Background(), account.ID)

	require.NoError(t, err)

	response, err := queries.GetAccount(context.Background(), account.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, response)
}

func TestGetAccount(t *testing.T) {

	account := createRandomAccount(t)

	response, err := queries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)

	require.NotEmpty(t, response)

	require.Equal(t, account.ID, response.ID)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	args := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := queries.ListAccounts(context.Background(), args)

	require.NoError(t, err)

	require.NotEmpty(t, accounts)

	require.Equal(t, len(accounts), 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestUpdateAccount(t *testing.T) {

	account := createRandomAccount(t)

	args := UpdateAccountParams{
		ID:      account.ID,
		Balance: 50,
	}

	response, err := queries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)

	require.NotEmpty(t, account)

	require.Equal(t, account.ID, response.ID)
	require.Equal(t, response.Owner, account.Owner)
	require.Equal(t, response.Currency, account.Currency)
	require.Equal(t, args.Balance, response.Balance)
}
