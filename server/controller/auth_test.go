package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	mockdb "github.com/gitaepark/carrot-market/db/mock"
	db "github.com/gitaepark/carrot-market/db/sqlc"
	"github.com/gitaepark/carrot-market/dto"
	"github.com/gitaepark/carrot-market/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	user, password := createRandomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":    user.Email,
				"password": password,
				"nickname": user.Nickname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireUserResponseMatch(t, recorder.Body, user)
			},
		},
		{
			name: "RequiredEmail",
			body: gin.H{
				"password": password,
				"nickname": user.Nickname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireErrorMatch(t, recorder.Body, util.ErrRequired("email"))
			},
		},
		{
			name: "InvalidEmailType",
			body: gin.H{
				"email":    1,
				"password": password,
				"nickname": user.Nickname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireErrorMatch(t, recorder.Body, util.ErrType("email", "string"))
			},
		},
		{
			name: "InvalidEmailFormat",
			body: gin.H{
				"email":    util.CreateRandomString(10),
				"password": password,
				"nickname": user.Nickname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireErrorMatch(t, recorder.Body, util.ErrEmail("email"))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			controller := newTestController(t, store)

			tc.buildStubs(store)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/auth/register"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			controller.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func createRandomUser(t *testing.T) (db.User, string) {
	password := util.CreateRandomPassword()
	hashedPassword, _ := util.HashPassword(password)

	user := db.User{
		UserID:         util.CreateRandomInt32(1, 30),
		Email:          util.CreateRandomEmail(),
		HashedPassword: hashedPassword,
		Nickname:       util.CreateRandomNickname(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return user, password
}

func requireUserResponseMatch(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var rsp dto.UserResponse
	err = json.Unmarshal(data, &rsp)

	require.NoError(t, err)
	require.Equal(t, user.Email, rsp.Email)
	require.Equal(t, user.Email, rsp.Email)
	require.Equal(t, user.Nickname, rsp.Nickname)
	require.WithinDuration(t, user.CreatedAt, rsp.CreatedAt, time.Second)
	require.WithinDuration(t, user.UpdatedAt, rsp.UpdatedAt, time.Second)
}
