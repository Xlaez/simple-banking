package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	arg := CreateEntryParams {
		Amount: account.Balance,
		AccountID: account.ID,
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	
	require.NotEmpty(t, entry.CreatedAt)
	require.NotEmpty(t, entry.ID)

	return entry;
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)

	entry, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	require.Equal(t, entry1.Amount, entry.Amount)
	require.Equal(t, entry1.AccountID, entry.AccountID)
}

func TestUpdateEntry(t *testing.T) {
	entry1 := createRandomEntry(t)

	arg := UpdateEntryParams{
		ID: entry1.ID,
		Amount: entry1.Amount,
	}

	entry, err := testQueries.UpdateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, entry.ID)

	require.Equal(t, entry.Amount, entry1.Amount)
	require.Equal(t, entry.AccountID, entry1.AccountID)

}

func TestListEntries(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams {
		Limit: 5,
		Offset: 5,
	}

	entry, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entry, 5)

	for _, list := range entry {
		require.NotEmpty(t, list)
	}

}

func TestDeleteEntry(t *testing.T) {
	entry := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)

	require.NoError(t, err)

	newEntry, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, newEntry)
}