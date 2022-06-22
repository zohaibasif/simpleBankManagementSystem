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

func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	args := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomAmount(),
	}

	transfer, err := queries.CreateTransfer(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, args.ToAccountID, transfer.ToAccountID)
	require.Equal(t, args.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1, account2 := createRandomAccount(t), createRandomAccount(t)

	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1, account2 := createRandomAccount(t), createRandomAccount(t)

	transfer := createRandomTransfer(t, account1, account2)

	response, err := queries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, response)
	require.Equal(t, transfer.ID, response.ID)
}

func TestListTransfers(t *testing.T) {

	account1, account2 := createRandomAccount(t), createRandomAccount(t)

	for i := 0; i < 10; i++ {
		queries.CreateTransfer(context.Background(), CreateTransferParams{
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount:        util.RandomAmount(),
		})
	}

	args := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := queries.ListTransfers(context.Background(), args)

	require.NoError(t, err)

	require.NotEmpty(t, transfers)

	require.Equal(t, len(transfers), 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
