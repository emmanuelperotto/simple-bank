package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"simplebank/util"
	"testing"
)

func TestQueries_CreateEntry(t *testing.T) {
	tests := []struct {
		name   string
		params CreateEntryParams
	}{
		{
			name: "When successfully creates it with positive amount",
			params: CreateEntryParams{
				Amount: 200,
			},
		},
		{
			name: "When successfully creates it with negative amount",
			params: CreateEntryParams{
				Amount: -150,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			account, err := randomAccount(ctx)
			require.NoError(t, err)
			require.NotZero(t, account.ID)

			entry, err := testQueries.CreateEntry(ctx, CreateEntryParams{
				AccountID: account.ID,
				Amount:    tt.params.Amount,
			})

			require.Equal(t, tt.params.Amount, entry.Amount)
			require.Equal(t, account.ID, entry.AccountID)

			require.NotZero(t, entry.ID)
			require.NotZero(t, entry.CreatedAt)
		})
	}
}

func TestQueries_GetEntry(t *testing.T) {
	tests := []struct {
		name        string
		testingFunc func(t *testing.T)
	}{
		{
			name: "When entry exists",
			testingFunc: func(t *testing.T) {
				ctx := context.Background()

				account, err := randomAccount(ctx)
				require.NoError(t, err)
				require.NotZero(t, account.ID)

				entry, err := testQueries.CreateEntry(ctx, CreateEntryParams{
					AccountID: account.ID,
					Amount:    util.RandomMoney(),
				})
				require.NoError(t, err)
				require.NotZero(t, entry.ID)

				e, err := testQueries.GetEntry(ctx, entry.ID)
				require.NoError(t, err)
				require.Equal(t, entry.ID, e.ID)
			},
		},
		{
			name: "When entry does not exist",
			testingFunc: func(t *testing.T) {
				ctx := context.Background()

				_, err := testQueries.GetEntry(ctx, 0)
				require.Error(t, err)
				require.Equal(t, sql.ErrNoRows, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.testingFunc)
	}
}

func TestQueries_ListEntries(t *testing.T) {
	tests := []struct {
		name        string
		testingFunc func(t *testing.T)
	}{
		{
			name: "When it has enough entries to show",
			testingFunc: func(t *testing.T) {
				ctx := context.Background()
				numEntries := 5

				account, err := randomAccount(ctx)
				require.NoError(t, err)
				require.NotZero(t, account.ID)

				for i := 0; i < numEntries; i++ {
					entry, err := testQueries.CreateEntry(ctx, CreateEntryParams{
						AccountID: account.ID,
						Amount:    util.RandomMoney(),
					})
					require.NoError(t, err)
					require.NotZero(t, entry.ID)
				}

				entries, err := testQueries.ListEntries(ctx, ListEntriesParams{
					Limit:  int32(numEntries),
					Offset: 0,
				})

				require.NoError(t, err)
				require.Len(t, entries, numEntries)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.testingFunc)
	}
}
