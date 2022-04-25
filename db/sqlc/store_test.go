package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	t.Parallel()
	store := NewStore(testDb)
	ctx := context.Background()

	fromAcc, err := randomAccount(ctx)
	require.NoError(t, err)

	toAcc, err := randomAccount(ctx)
	require.NoError(t, err)

	//Run n concurrent transfer transactions
	n := 5
	errors := make(chan error)
	results := make(chan TransferTxResult)
	amount := int64(10)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAcc.ID,
				ToAccountID:   toAcc.ID,
				Amount:        amount,
			})

			errors <- err
			results <- result
		}()
	}

	//Assert results
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//transfer assertions
		transfer := result.Transfer
		require.Equal(t, fromAcc.ID, transfer.FromAccountID)
		require.Equal(t, toAcc.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		//entries assertions
		fromEntry := result.FromEntry
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		require.Equal(t, -amount, fromEntry.Amount)
		require.Equal(t, fromAcc.ID, fromEntry.AccountID)

		toEntry := result.ToEntry
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		require.Equal(t, amount, toEntry.Amount)
		require.Equal(t, toAcc.ID, toEntry.AccountID)

		//accounts assertions
		fromAccResult := result.FromAccount
		toAccResult := result.ToAccount
		require.Equal(t, fromAcc.Balance+fromEntry.Amount, fromAccResult.Balance)
		require.Equal(t, toAcc.Balance+toEntry.Amount, toAccResult.Balance)

		//Update in-memory accounts for the next assertion
		fromAcc = fromAccResult
		toAcc = toAccResult
	}
}

func TestStore_TransferTxDeadLock(t *testing.T) {
	t.Parallel()
	store := NewStore(testDb)
	ctx := context.Background()

	fromAcc, err := randomAccount(ctx)
	require.NoError(t, err)

	toAcc, err := randomAccount(ctx)
	require.NoError(t, err)

	//Run n concurrent transfer transactions
	n := 10
	errors := make(chan error)
	amount := int64(10)
	for i := 0; i < n; i++ {
		fromAccID := fromAcc.ID
		toAccID := toAcc.ID

		if i%2 == 0 {
			fromAccID = toAcc.ID
			toAccID = fromAcc.ID
		}

		go func() {
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccID,
				ToAccountID:   toAccID,
				Amount:        amount,
			})

			errors <- err
		}()
	}

	//Assert no error
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)
	}

	//Assert for final balance
	account, err := store.GetAccount(ctx, fromAcc.ID)
	require.NoError(t, err)

	require.Equal(t, fromAcc.Balance, account.Balance)

	a, err := store.GetAccount(ctx, toAcc.ID)
	require.NoError(t, err)

	require.Equal(t, toAcc.Balance, a.Balance)
}
