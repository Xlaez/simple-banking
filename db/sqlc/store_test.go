package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T){
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 5

	amount := int64(10) // amount to be transferred

	// build a channel to communicate with threads: this channel returns result and error to main routine to validate
	errs := make(chan error)
	results := make(chan TransferTxResult)


	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})

			// send result and err to channels
			errs <- err
			results <- result
		}()
	}
	// Note: to send to a channel the chanel is on the left whereas to receive, we place the channel on the right
	// check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)
		
		result := <- results
		fmt.Println("results:", result.FromAccount.ID, account1.ID)
		require.NotEmpty(t, result)

		// check transfer of result object

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries of result object

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)
		
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t,toAccount.ID, account2.ID)
		
		// check account balance of result object
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check final balance

	updatedAcc1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAcc2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println("After =======> ", updatedAcc1.Balance, updatedAcc2.Balance)
	require.Equal(t, int(account1.Balance) -n * int(amount) , int(updatedAcc1.Balance))
	require.Equal(t, int(account2.Balance) +n * int(amount), int(updatedAcc2.Balance))
}