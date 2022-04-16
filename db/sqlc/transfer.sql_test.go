package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"simplebank/util"
	"testing"
)

func randomTransfer(ctx context.Context) (Transfer, error) {
	fromAcc, err := randomAccount(ctx)
	if err != nil {
		return Transfer{}, err
	}

	toAcc, err := randomAccount(ctx)
	if err != nil {
		return Transfer{}, err
	}

	return testQueries.CreateTransfer(ctx, CreateTransferParams{
		FromAccountID: fromAcc.ID,
		ToAccountID:   toAcc.ID,
		Amount:        util.RandomMoney(),
	})
}

func TestQueries_CreateTransfer(t *testing.T) {
	tests := []struct {
		name        string
		testingFunc func(t *testing.T)
	}{
		{
			name: "When it successfully creates it",
			testingFunc: func(t *testing.T) {
				ctx := context.Background()

				fromAcc, err := randomAccount(ctx)
				require.NoError(t, err)
				require.NotZero(t, fromAcc.ID)

				toAcc, err := randomAccount(ctx)
				require.NoError(t, err)
				require.NotZero(t, toAcc.ID)

				amount := util.RandomMoney()
				transfer, err := testQueries.CreateTransfer(ctx, CreateTransferParams{
					FromAccountID: fromAcc.ID,
					ToAccountID:   toAcc.ID,
					Amount:        amount,
				})
				require.NoError(t, err)
				require.NotZero(t, transfer.ID)
				require.NotEmpty(t, transfer.CreatedAt)

				require.Equal(t, fromAcc.ID, transfer.FromAccountID)
				require.Equal(t, toAcc.ID, transfer.ToAccountID)
				require.Equal(t, amount, transfer.Amount)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.testingFunc)
	}
}

func TestQueries_GetTransfer(t *testing.T) {

	tests := []struct {
		name        string
		testingFunc func(t *testing.T)
	}{
		{
			name: "When transfer exists",
			testingFunc: func(t *testing.T) {
				ctx := context.Background()
				transfer, err := randomTransfer(ctx)
				require.NoError(t, err)
				require.NotZero(t, transfer.ID)

				tr, err := testQueries.GetTransfer(ctx, transfer.ID)
				require.NoError(t, err)
				require.Equal(t, transfer.ID, tr.ID)
			},
		},
		{
			name: "When transfer does not exist",
			testingFunc: func(t *testing.T) {
				ctx := context.Background()

				_, err := testQueries.GetTransfer(ctx, 0)
				require.Error(t, err)
				require.Equal(t, sql.ErrNoRows, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.testingFunc)
	}
}

func TestQueries_ListTransfers(t *testing.T) {

	tests := []struct {
		name        string
		testingFunc func(t *testing.T)
	}{
		{
			name: "When there is enough transfers to list",
			testingFunc: func(t *testing.T) {
				ctx := context.Background()
				numTransfers := 5

				for i := 0; i < numTransfers; i++ {
					transfer, err := randomTransfer(ctx)
					require.NoError(t, err)
					require.NotZero(t, transfer.ID)
				}

				transfers, err := testQueries.ListTransfers(ctx, ListTransfersParams{
					Limit:  int32(numTransfers),
					Offset: 0,
				})
				require.NoError(t, err)
				require.Len(t, transfers, numTransfers)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.testingFunc)
	}
}
