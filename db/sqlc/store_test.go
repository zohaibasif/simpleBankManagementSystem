package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	args := TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        11,
	}

	errors := make(chan error)
	results := make(chan TransferTxResult)
	existed := make(map[int]bool)

	n := 5

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), args)
			errors <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, args.Amount)
		require.NotZero(t, transfer.CreatedAt)
		require.NotZero(t, transfer.ID)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check from account entry
		fromEntry := result.FromEntry

		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -args.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//check to account entry
		toEntry := result.ToEntry

		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, args.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)
		require.NotZero(t, fromAccount.ID)
		require.NotZero(t, fromAccount.CreatedAt)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)
		require.NotZero(t, toAccount.ID)
		require.NotZero(t, toAccount.CreatedAt)

		dif1 := account1.Balance - fromAccount.Balance
		dif2 := toAccount.Balance - account2.Balance

		require.Equal(t, dif1, dif2)                     // same amount of money got out of account 1 and added to account 2
		require.True(t, dif1 > 0)                        // as there cannot be a negative transfer
		require.True(t, int(dif1)%int(args.Amount) == 0) // if there are 5 transactions from account 1 of amount 10 to account 2, then 50%10 should be equal to zero

		k := int(dif1 / args.Amount)       // no of transactions
		require.True(t, k > 0 && k <= n)   // k should be under 1 to n
		require.NotContains(t, existed, k) // k should be unique each time
		existed[k] = true
	}

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-(int64(n)*args.Amount), updatedAccount1.Balance)
	require.Equal(t, account2.Balance+(int64(n)*args.Amount), updatedAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	errors := make(chan error)

	n := 10

	for i := 0; i < n; i++ {

		fromAccount := account1.ID
		toAccount := account2.ID

		if i%2 == 1 {
			fromAccount = account2.ID
			toAccount = account1.ID
		}

		args := TransferTxParams{
			FromAccountID: fromAccount,
			ToAccountID:   toAccount,
			Amount:        11,
		}

		go func() {
			_, err := store.TransferTx(context.Background(), args)
			errors <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)
	}

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

}
