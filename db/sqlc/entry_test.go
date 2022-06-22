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

func createRandomEntry(t *testing.T, account Account) Entry {
	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomAmount(),
	}

	entry, err := queries.CreateEntry(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, account)
	response, err := queries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, response)
	require.Equal(t, response.ID, entry.ID)
}

func TestListEntries(t *testing.T) {

	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		queries.CreateEntry(context.Background(), CreateEntryParams{
			AccountID: account.ID,
			Amount:    util.RandomAmount(),
		})
	}

	args := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := queries.ListEntries(context.Background(), args)

	require.NoError(t, err)

	require.NotEmpty(t, entries)

	require.Equal(t, len(entries), 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
