package db

import (
	"context"
	"database/sql"
	"log"
)

var txKey = struct{}{}

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SqlStore struct {
	*Queries
	db *sql.DB
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Tranfers    Transfer
	FromAccount Account
	ToAccount   Account
	FromEntry   Entry
	ToEntry     Entry
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SqlStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SqlStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Fatalf("tx error: %v, rb error: %v", err, rbErr)
			return rbErr
		}
		return err
	}

	return tx.Commit()
}

func (store *SqlStore) TransferTx(ctx context.Context, params TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Tranfers, err = q.CreateTranfer(ctx, CreateTranferParams{
			FromAcountID: params.FromAccountID,
			ToAcountID:   params.ToAccountID,
			Amount:       params.Amount,
		})

		if err != nil {
			return err
		}

		result.FromAccount, result.ToAccount, err = store.transferBalance(ctx, params.FromAccountID, params.ToAccountID, params.Amount)
		if err != nil {
			return err
		}

		result.FromEntry, result.ToEntry, err = store.craeteTransferEntires(ctx, params.FromAccountID, params.ToAccountID, params.Amount)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func (store *SqlStore) transferBalance(ctx context.Context, fromAccountID int64, toAccountID int64, amount int64) (from Account, to Account, err error) {
	if fromAccountID < toAccountID {
		from, err = store.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     fromAccountID,
			Amount: -amount,
		})
		if err != nil {
			return
		}

		to, err = store.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     toAccountID,
			Amount: amount,
		})

	} else {
		to, err = store.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     toAccountID,
			Amount: amount,
		})

		if err != nil {
			return
		}

		from, err = store.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     fromAccountID,
			Amount: -amount,
		})
	}

	return
}

func (store *SqlStore) craeteTransferEntires(ctx context.Context, fromAccountID, toAccountID, amount int64) (from Entry, to Entry, err error) {
	from, err = store.CreateEntry(ctx, CreateEntryParams{
		OwnerID: fromAccountID,
		Amount:  -amount,
	})
	if err != nil {
		return
	}

	to, err = store.CreateEntry(ctx, CreateEntryParams{
		OwnerID: toAccountID,
		Amount:  amount,
	})

	return
}
