package db

import (
	"context"
	"database/sql"
	"testing"

	"simple-bank/util"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer{
	transfer1 := createRandomAccount(t)
	transfer2 := createRandomAccount(t)

	arg := CreateTransferParams {
		Amount: util.RandomMoney(),
		FromAccountID: transfer1.ID,
		ToAccountID: transfer2.ID,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.NotEmpty(t, transfer.CreatedAt)
	require.NotEmpty(t, transfer.ID)

	return transfer;
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	randomTransfer := createRandomTransfer(t)

	transfer, err := testQueries.GetTransfer(context.Background(), randomTransfer.ID)

	require.NoError(t, err)
	require.NotZero(t, transfer.ID)

	require.Equal(t, transfer.Amount, randomTransfer.Amount)
	require.Equal(t, transfer.FromAccountID, randomTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, randomTransfer.ToAccountID)

}

func TestUpdateTransfer(t *testing.T) {
	randomTransfer := createRandomTransfer(t)

	arg := UpdateTransferParams {
		ID: randomTransfer.ID,
		Amount: util.RandomMoney(),
	}

	transfer, err := testQueries.UpdateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotZero(t, transfer.CreatedAt)

	require.Equal(t, transfer.FromAccountID, randomTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, randomTransfer.ToAccountID)
}

func TestListTransfer(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	arg := ListTransfersParams {
		Limit: 5,
		Offset: 5,
	}

	transfer, err  := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfer, 5)

	for _, list := range transfer {
		require.NotEmpty(t, list)
	} 

}

func TestDeleteTransfer(t *testing.T) {
	randomTransfer := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), randomTransfer.ID)

	require.NoError(t, err)

	transfer, err := testQueries.GetTransfer(context.Background(), randomTransfer.ID)
	
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer)
}