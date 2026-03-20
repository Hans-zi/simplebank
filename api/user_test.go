package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/Hans-zi/simple_bank/db/mock"
	db "github.com/Hans-zi/simple_bank/db/sqlc"
	"github.com/Hans-zi/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type eqCreateUserMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserMatcher) Matches(x any) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(arg.HashedPassword, e.password)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %s", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserMatcher{arg, password}
}

func TestCreateUserApi(t *testing.T) {
	user, password := createRandomUser(t)
	testcases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}

				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(args, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireMatchBodyUser(t, user, password, recorder.Body)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}

				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(args, password)).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username":  user.Username,
				"password":  "123",
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := fmt.Sprintf("/users")

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, recorder)
		})
	}
}

func requireMatchBodyUser(t *testing.T, user db.User, password string, body *bytes.Buffer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var res userResponse
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	require.Equal(t, res.Username, user.Username)
	require.Equal(t, res.FullName, user.FullName)
	require.Equal(t, res.Email, user.Email)
}

func createRandomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	return db.User{
		Username:       util.RandomString(6),
		HashedPassword: hashedPassword,
		FullName:       util.RandomString(6),
		Email:          util.RandomEmail(),
	}, password
}
