package db

import (
	"context"
	"testing"

	"github.com/mosmatan/simplebank/utils"
)

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    utils.RandomName(),
		Balance:  utils.RandomBalance(),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if account.Owner != arg.Owner {
		t.Error("owner mismatch")
	}
	if account.Balance != arg.Balance {
		t.Error("balance mismatch")
	}
	if account.Currency != arg.Currency {
		t.Error("currency mismatch")
	}
	if account.ID == 0 {
		t.Error("missing id")
	}
}

func TestGetAccount(t *testing.T) {
	createAccount := CreateAccountParams{
		Owner:    utils.RandomName(),
		Balance:  utils.RandomBalance(),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), createAccount)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if account2.Owner != account.Owner {
		t.Error("owner mismatch")
	}
	if account2.Balance != account.Balance {
		t.Error("balance mismatch")
	}
	if account2.Currency != account.Currency {
		t.Error("currency mismatch")
	}
	if account2.ID != account.ID {
		t.Error("id mismatch")
	}
}

func TestUpdateAccount(t *testing.T) {
	createAccount := CreateAccountParams{
		Owner:    utils.RandomName(),
		Balance:  utils.RandomBalance(),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), createAccount)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	update := UpdateAccountParams{
		ID:      account.ID,
		Balance: 300,
	}
	account2, err := testQueries.UpdateAccount(context.Background(), update)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if account2.Owner != account.Owner {
		t.Error("owner mismatch")
	}
	if account2.Balance != update.Balance {
		t.Error("balance mismatch")
	}
	if account2.Currency != account.Currency {
		t.Error("currency mismatch")
	}
	if account2.ID != account.ID {
		t.Error("id mismatch")
	}
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		arg := CreateAccountParams{
			Owner:    utils.RandomName(),
			Balance:  utils.RandomBalance(),
			Currency: "USD",
		}
		_, err := testQueries.CreateAccount(context.Background(), arg)
		if err != nil {
			t.Errorf("error: %v", err)
		}
	}
	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
		Limit:  5,
		Offset: 5,
	})
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if len(accounts) != 5 {
		t.Errorf("expected 5 accounts, got %d", len(accounts))
	}
	for _, account := range accounts {
		if account.Currency != "USD" {
			t.Error("currency mismatch should be USD but got", account.Currency)
		}
	}
}

func TestDeleteAccount(t *testing.T) {
	createAccount := CreateAccountParams{
		Owner:    utils.RandomName(),
		Balance:  utils.RandomBalance(),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), createAccount)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	err = testQueries.DeleteAccount(context.Background(), account.ID)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	if err == nil {
		t.Errorf("expected error, but got %v", account2)
	}
}
