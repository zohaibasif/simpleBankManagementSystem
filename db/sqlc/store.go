package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	Db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return Store{
		Db:      db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {

	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		if rbkErr := tx.Rollback(); err != nil {
			return fmt.Errorf("tx error %v, rollback error %v", err, rbkErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer
	FromAccount Account
	ToAccount   Account
	FromEntry   Entry
	ToEntry     Entry
}

func (s *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	result := TransferTxResult{}

	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		// update accounts' balance
		if args.FromAccountID > args.ToAccountID {
			result.FromAccount, result.ToAccount, err = addBalance(ctx, q, args.FromAccountID, -args.Amount, args.ToAccountID, args.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addBalance(ctx, q, args.ToAccountID, args.Amount, args.FromAccountID, -args.Amount)
		}
		return nil
	})

	return result, err
}

func addBalance(ctx context.Context, q *Queries, account1Id, amount1, account2Id, amount2 int64) (account1 Account, account2 Account, err error) {

	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account1Id,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account2Id,
		Amount: amount2,
	})

	return
}
