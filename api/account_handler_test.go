package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	db "simplebank/db/sqlc"
	mockdb "simplebank/db/sqlc/mock"
	"simplebank/util"
	"testing"
	"time"
)

var (
	defaultCreatedAt = time.Date(2022, time.April, 24, 21, 18, 0, 0, time.UTC)
)

type stub struct {
	store db.Store
}

func Test_accountHandler_get(t *testing.T) {
	tests := []struct {
		name          string
		accountID     int64
		buildStubs    func(ctrl *gomock.Controller) stub
		runAssertions func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "When it successfully finds the account",
			accountID: 10,
			buildStubs: func(ctrl *gomock.Controller) stub {
				store := mockdb.NewMockStore(ctrl)
				account := db.Account{
					ID:        10,
					Owner:     "Perotto",
					Balance:   100,
					Currency:  "USD",
					CreatedAt: defaultCreatedAt,
				}

				store.EXPECT().GetAccount(gomock.Any(), int64(10)).
					Times(1).
					Return(account, nil)

				return stub{store: store}
			},
			runAssertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var responseBody db.Account
				require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))

				wantResponseBody := db.Account{
					ID:        10,
					Owner:     "Perotto",
					Balance:   100,
					Currency:  "USD",
					CreatedAt: defaultCreatedAt,
				}

				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Equal(t, wantResponseBody, responseBody)
			},
		},

		{
			name:      "When sending id less than 1",
			accountID: 0,
			buildStubs: func(ctrl *gomock.Controller) stub {
				store := mockdb.NewMockStore(ctrl)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)

				return stub{store: store}
			},
			runAssertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var responseBody gin.H
				require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))

				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				assert.Equal(t, "Key: 'getAccountRequest.ID' Error:Field validation for 'ID' failed on the 'required' tag", responseBody["error"])
			},
		},
		{
			name:      "When account not found",
			accountID: util.RandomInt(1, 5),
			buildStubs: func(ctrl *gomock.Controller) stub {
				store := mockdb.NewMockStore(ctrl)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)

				return stub{store: store}
			},
			runAssertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var responseBody gin.H
				require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))

				assert.Equal(t, http.StatusNotFound, recorder.Code)
				assert.Equal(t, "account not found", responseBody["error"])
			},
		},
		{
			name:      "When there is a generic error fetching account",
			accountID: 1,
			buildStubs: func(ctrl *gomock.Controller) stub {
				store := mockdb.NewMockStore(ctrl)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, errors.New("run, it's all broken"))

				return stub{store: store}
			},
			runAssertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var responseBody gin.H
				require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))

				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
				assert.Equal(t, "unknown error", responseBody["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			//Builds stubs
			stubs := tt.buildStubs(ctrl)

			//Start test server and send request
			url := fmt.Sprintf("/accounts/%d", tt.accountID)
			server := NewServer(stubs.store)
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			//Assertions
			tt.runAssertions(t, recorder)
		})
	}
}

func Test_accountHandler_post(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   createAccountRequest
		buildStubs    func(ctrl *gomock.Controller) stub
		runAssertions func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "When it succeeds",
			requestBody: createAccountRequest{
				Owner:    "Emmanuel Perotto",
				Currency: "USD",
			},
			buildStubs: func(ctrl *gomock.Controller) stub {
				store := mockdb.NewMockStore(ctrl)
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).
					Return(db.Account{
						ID:        1,
						Owner:     "Emmanuel Perotto",
						Balance:   0,
						Currency:  "USD",
						CreatedAt: defaultCreatedAt,
					}, nil)

				return stub{
					store: store,
				}
			},
			runAssertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var responseBody gin.H

				require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
				assert.Equal(t, http.StatusCreated, recorder.Code)

				wantResponseBody := gin.H{
					"id":         float64(1),
					"owner":      "Emmanuel Perotto",
					"currency":   "USD",
					"balance":    float64(0),
					"created_at": "2022-04-24T21:18:00Z",
				}

				assert.Equal(t, wantResponseBody, responseBody)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			//Builds stubs
			stubs := tt.buildStubs(ctrl)

			//Start test server and send request
			url := "/accounts"
			server := NewServer(stubs.store)
			recorder := httptest.NewRecorder()

			bodyBytes, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			//Assertions
			tt.runAssertions(t, recorder)
		})
	}
}
