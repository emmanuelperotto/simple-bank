package db

import (
	"context"
	"database/sql"
	"fmt"
)

//Store provides all functions to execute db queries and transactions
type (
	Store struct {
		*Queries
		db *sql.DB
	}
	//TransferTxParams contains the input parameters of the transfer transaction
	TransferTxParams struct {
		FromAccountID int64 `json:"from_account_id"`
		ToAccountID   int64 `json:"to_account_id"`
		Amount        int64 `json:"amount"`
	}
	//TransferTxResult is the result of the transfer transaction
	TransferTxResult struct {
		Transfer    `json:"transfer"`
		FromAccount Account `json:"from_account"`
		ToAccount   Account `json:"to_account"`
		FromEntry   Entry   `json:"from_entry"`
		ToEntry     Entry   `json:"to_entry"`
	}

	updateBalanceRequest struct {
		fromId     int64
		fromAmount int64
		toId       int64
		toAmount   int64
	}
)

//NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (s *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %+v, rb err: %+v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

//TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries and update accounts' balance within a single database transaction
func (s *Store) TransferTx(ctx context.Context, params TransferTxParams) (result TransferTxResult, err error) {
	err = s.execTx(ctx, func(queries *Queries) error {
		if result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: params.FromAccountID,
			ToAccountID:   params.ToAccountID,
			Amount:        params.Amount,
		}); err != nil {
			return err
		}

		if result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.FromAccountID,
			Amount:    -params.Amount,
		}); err != nil {
			return err
		}

		if result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.ToAccountID,
			Amount:    params.Amount,
		}); err != nil {
			return err
		}

		result.FromAccount, result.ToAccount, err = updateAccountBalances(ctx, queries, updateBalanceRequest{
			fromId:     params.FromAccountID,
			fromAmount: result.FromEntry.Amount,
			toId:       params.ToAccountID,
			toAmount:   result.ToEntry.Amount,
		})

		return err
	})

	return result, err
}

//updateAccountBalances updates account balances ensuring the smallest id will be updated first to avoid deadlocks
func updateAccountBalances(ctx context.Context, q *Queries, request updateBalanceRequest) (fromAcc Account, toAcc Account, err error) {
	fromId := request.fromId
	toId := request.toId
	fromAmount := request.fromAmount
	toAmount := request.toAmount

	if fromId < toId {
		fromAcc, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			Amount: fromAmount,
			ID:     fromId,
		})

		toAcc, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			Amount: toAmount,
			ID:     toId,
		})

		return
	}
	toAcc, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: toAmount,
		ID:     toId,
	})

	fromAcc, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: fromAmount,
		ID:     fromId,
	})

	return
}
