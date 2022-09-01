package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"simple-bank/util"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T)  Account{
	// define a test data request
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner: user.Username,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	// handle test validation

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	// check to make sure accountId is not zero

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T){
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	createAcc := createRandomAccount(t)
	getAcc, err := testQueries.GetAccount(context.Background(), createAcc.ID)

	require.NoError(t, err)
	require.NotEmpty(t, getAcc)

	require.Equal(t, createAcc.ID, getAcc.ID)	
	require.Equal(t, createAcc.Owner, getAcc.Owner)
	require.Equal(t, createAcc.Balance, getAcc.Balance)
	require.Equal(t, createAcc.Currency, getAcc.Currency)	
	require.Equal(t, createAcc.CreatedAt, getAcc.CreatedAt)	
}

func TestUpdateAccount(t *testing.T) {
	oldAcc := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID: oldAcc.ID,
		Balance: util.RandomMoney(),
	}

	newAcc, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, newAcc)
	
	require.Equal(t, oldAcc.Owner, newAcc.Owner)
	require.Equal(t, oldAcc.Currency, newAcc.Currency)	

	require.WithinDuration(t, oldAcc.CreatedAt, newAcc.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t) 

	err := testQueries.DeleteAccount(context.Background(), account.ID)

	require.NoError(t, err)


	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}