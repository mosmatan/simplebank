package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/mosmatan/simplebank/db/mock"
	db "github.com/mosmatan/simplebank/db/sqlc"
	"github.com/mosmatan/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := createRandomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	request := httptest.NewRequest("GET", url, nil)

	server.router.ServeHTTP(recorder, request)

	require.Equal(t, recorder.Code, 200)
	requireAccountEqual(t, recorder.Body.Bytes(), account)
}

func createRandomAccount() db.Account {
	return db.Account{
		ID:       int64(utils.RandomInt(1, 1000)),
		Owner:    utils.RandomName(),
		Balance:  utils.RandomBalance(),
		Currency: "USD",
	}
}

func requireAccountEqual(t *testing.T, body []byte, account db.Account) {
	var gotAccount db.Account
	err := json.Unmarshal(body, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
