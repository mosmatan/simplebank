package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/mosmatan/simplebank/utils"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

// TestMain function is used to setup the test environment
func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open("postgres", "postgresql://root:secret@localhost:5432/bank?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomName(),
		Balance:  utils.RandomBalance(),
		Currency: "USD",
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	return account
}
