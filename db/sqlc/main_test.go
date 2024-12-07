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
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
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
