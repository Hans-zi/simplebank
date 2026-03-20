package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/Hans-zi/simple_bank/db/mock"
	db "github.com/Hans-zi/simple_bank/db/sqlc"
	"github.com/Hans-zi/simple_bank/token"
	"github.com/Hans-zi/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTransferApi(t *testing.T) {
	user1, _ := createRandomUser(t)
	user2, _ := createRandomUser(t)
	user3, _ := createRandomUser(t)

	account1 := createRandomAccount(user1.Username)
	account2 := createRandomAccount(user2.Username)
	account3 := createRandomAccount(user3.Username)

	account1.Currency = "USD"
	account2.Currency = "USD"
	account3.Currency = "EUR"

	var amount int64 = 10
	transfer := createTransfer(account1.ID, account2.ID, amount)
	fromEntry := createEntry(account1.ID, -amount)
	toEntry := createEntry(account2.ID, amount)
	transferTXResult := db.TransferTxResult{
		Transfer:    transfer,
		FromAccount: account1,
		ToAccount:   account2,
		FromEntry:   fromEntry,
		ToEntry:     toEntry,
	}

	testcases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationHeaderTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)
				txArgs := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}

				store.EXPECT().TransferTX(gomock.Any(), gomock.Eq(txArgs)).
					Times(1).
					Return(transferTXResult, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireMatchBodyTransferTx(t, transferTXResult, recorder.Body)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server := NewTestServer(t, store)

			recorder := httptest.NewRecorder()
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/transfer")
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, recorder)
		})
	}
}

func requireMatchBodyTransferTx(t *testing.T, result db.TransferTxResult, body *bytes.Buffer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var res db.TransferTxResult
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)

	require.NotEmpty(t, res)
	require.Equal(t, result.Transfer, res.Transfer)
	require.Equal(t, result.FromAccount, res.FromAccount)
	require.Equal(t, result.ToAccount, res.ToAccount)
	require.Equal(t, result.FromEntry, res.FromEntry)
	require.Equal(t, result.ToEntry, res.ToEntry)

}

func createTransfer(from, to, amount int64) db.Transfer {
	return db.Transfer{
		ID:            util.RandomInt(1, 1000),
		FromAccountID: from,
		ToAccountID:   to,
		Amount:        amount,
	}
}

func createEntry(id, amount int64) db.Entry {
	return db.Entry{
		ID:        util.RandomInt(1, 1000),
		AccountID: id,
		Amount:    amount,
	}
}
