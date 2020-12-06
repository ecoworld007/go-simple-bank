package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("before:::::::::", account1.Balance, account2.Balance)
	// run n concurrent transfer
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				Amount:        amount,
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
			})
			errs <- err
			results <- result
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.Amount, amount)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.NotZero(t, transfer.CreatedAt)
		require.NotZero(t, transfer.ID)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		// check ToEntry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)

		// check fromEntry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)

		// check balance update
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)

		fmt.Println("tx::::::::", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)
		k := int(diff1 / amount)
		require.True(t, k >= 1 || k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final update to accounts after all transfers
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
	fmt.Println("after:::::::::", updatedAccount1.Balance, updatedAccount2.Balance)
}
