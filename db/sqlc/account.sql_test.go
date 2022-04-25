package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"simplebank/util"
	"testing"
)

func TestQueries_CreateAccount(t *testing.T) {
	tests := []struct {
		name   string
		params CreateAccountParams
	}{
		{
			name: "When successfully creates it",
			params: CreateAccountParams{
				Owner:    util.RandomOwner(),
				Balance:  util.RandomMoney(),
				Currency: util.RandomCurrency(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			account, err := testQueries.CreateAccount(context.Background(), tt.params)
			require.NoError(t, err)
			require.NotEmpty(t, account)

			require.Equal(t, tt.params.Owner, account.Owner)
			require.Equal(t, tt.params.Balance, account.Balance)
			require.Equal(t, tt.params.Currency, account.Currency)

			require.NotZero(t, account.ID)
			require.NotZero(t, account.CreatedAt)
		})
	}
}

func TestQueries_GetAccount(t *testing.T) {
	tests := []struct {
		name        string
		testingFunc func(t *testing.T)
	}{
		{
			name: "When Account doesn't exist",
			testingFunc: func(t *testing.T) {
				acc, err := testQueries.GetAccount(context.Background(), -1)
				require.Error(t, err)
				require.Empty(t, acc)
				require.Equal(t, sql.ErrNoRows, err)
			},
		},
		{
			name: "When Account exists",
			testingFunc: func(t *testing.T) {
				ctx := context.Background()

				acc, err := randomAccount(ctx)
				require.NoError(t, err)
				require.NotEmpty(t, acc)

				a, err := testQueries.GetAccount(ctx, acc.ID)
				require.NoError(t, err)
				require.NotZero(t, a.ID)
				require.NotEmpty(t, a.CreatedAt)

				require.Equal(t, acc.Owner, a.Owner)
				require.Equal(t, acc.Balance, a.Balance)
				require.Equal(t, acc.Currency, a.Currency)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.testingFunc(t)
		})
	}
}

func TestQueries_DeleteAccount(t *testing.T) {
	tests := []struct {
		name        string
		testingFunc func(t *testing.T)
	}{
		{
			name: "When accounts exists",
			testingFunc: func(t *testing.T) {
				account, err := randomAccount(context.Background())

				require.NoError(t, err)

				err = testQueries.DeleteAccount(context.Background(), account.ID)
				require.NoError(t, err)
			},
		},

		{
			name: "When accounts doesn't exist",
			testingFunc: func(t *testing.T) {
				err := testQueries.DeleteAccount(context.Background(), -10)
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testingFunc(t)
		})
	}
}

func TestQueries_ListAccounts(t *testing.T) {
	tests := []struct {
		name        string
		testingFunc func(t *testing.T)
	}{
		{
			name: "When there are accounts to be shown",
			testingFunc: func(t *testing.T) {
				ctx := context.Background()
				numAccounts := 5

				for i := 0; i < numAccounts; i++ {
					account, err := randomAccount(ctx)

					require.NoError(t, err)
					require.NotZero(t, account.ID)
				}

				accounts, err := testQueries.ListAccounts(ctx, ListAccountsParams{
					Limit:  int32(numAccounts),
					Offset: 0,
				})
				require.NoError(t, err)
				require.Len(t, accounts, numAccounts)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testingFunc(t)
		})
	}
}

func TestQueries_UpdateAccount(t *testing.T) {
	tests := []struct {
		name        string
		testingFunc func(t *testing.T)
	}{
		{
			name: "Successful update",
			testingFunc: func(t *testing.T) {
				ctx := context.Background()

				account, err := randomAccount(ctx)
				require.NoError(t, err)
				require.NotZero(t, account.ID)

				newBalance := account.Balance + 10
				a, err := testQueries.UpdateAccount(ctx, UpdateAccountParams{
					Balance: newBalance,
					ID:      account.ID,
				})

				require.NoError(t, err)
				require.Equal(t, account.ID, a.ID)
				require.Equal(t, newBalance, a.Balance)
			},
		},

		{
			name: "When account doesn't exist",
			testingFunc: func(t *testing.T) {
				ctx := context.Background()

				_, err := testQueries.UpdateAccount(ctx, UpdateAccountParams{
					Balance: util.RandomMoney(),
					ID:      0,
				})

				require.Error(t, err)
				require.Equal(t, sql.ErrNoRows, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.testingFunc(t)
		})
	}
}

func randomAccount(ctx context.Context) (Account, error) {
	return testQueries.CreateAccount(ctx, CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	})
}
