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

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			if err != nil {
				errs <- err
			} else {
				results <- result
			}
		}()
	}

	for i := 0; i < n; i++ {
		select {
		case err := <-errs:
			require.NoError(t, err)
			fmt.Println("error", err)

		case result := <-results:
			require.NotEmpty(t, result)

			transfer := result.Tranfers
			require.NotEmpty(t, transfer)
			require.Equal(t, account1.ID, transfer.FromAcountID)
			require.Equal(t, account2.ID, transfer.ToAcountID)
			require.Equal(t, amount, transfer.Amount)
			require.NotZero(t, transfer.CreatedAt)
			require.NotZero(t, transfer.ID)

			_, err := store.GetTranfer(context.Background(), transfer.ID)
			require.NoError(t, err)

			fromEntry := result.FromEntry
			require.NotEmpty(t, fromEntry)
			require.Equal(t, account1.ID, fromEntry.OwnerID)
			require.Equal(t, -amount, fromEntry.Amount)
			require.NotZero(t, fromEntry.CreatedAt)
			require.NotZero(t, fromEntry.ID)

			_, err = store.GetEntry(context.Background(), fromEntry.ID)
			require.NoError(t, err)

			toEntry := result.ToEntry
			require.NotEmpty(t, toEntry)
			require.Equal(t, account2.ID, toEntry.OwnerID)
			require.Equal(t, amount, toEntry.Amount)
			require.NotZero(t, toEntry.CreatedAt)
			require.NotZero(t, toEntry.ID)

			_, err = store.GetEntry(context.Background(), toEntry.ID)
			require.NoError(t, err)

			fromAccount := result.FromAccount
			require.NotEmpty(t, fromAccount)
			require.Equal(t, account1.ID, fromAccount.ID)
			require.Equal(t, account1.Balance-amount, fromAccount.Balance)

			toAccount := result.ToAccount
			require.NotEmpty(t, toAccount)
			require.Equal(t, account2.ID, toAccount.ID)
			require.Equal(t, account2.Balance+amount, toAccount.Balance)

			account1.Balance = fromAccount.Balance
			account2.Balance = toAccount.Balance

			fmt.Println("from account ", fromAccount.Balance, "\nto account ", toAccount.Balance)
		}
	}
}

func TestTransferTxDaedLock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	balance1 := account1.Balance
	balance2 := account2.Balance

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			if err != nil {
				errs <- err
			} else {
				results <- result
			}
		}()

		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account2.ID,
				ToAccountID:   account1.ID,
				Amount:        amount,
			})

			if err != nil {
				errs <- err
			} else {
				results <- result
			}
		}()
	}

	for i := 0; i < n*2; i++ {
		select {
		case err := <-errs:
			require.NoError(t, err)
			fmt.Println("error", err)
		case result := <-results:
			require.NotEmpty(t, result)
		}
	}

	require.Equal(t, balance1, account1.Balance)
	require.Equal(t, balance2, account2.Balance)
}
