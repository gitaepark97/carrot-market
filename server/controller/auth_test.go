package controller

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/corpix/uarand"
	"github.com/gin-gonic/gin"
	mockdb "github.com/gitaepark/carrot-market/db/mock"
	db "github.com/gitaepark/carrot-market/db/sqlc"
	"github.com/gitaepark/carrot-market/dto"
	"github.com/gitaepark/carrot-market/token"
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
		{
			name: "RequiredPassword",
			body: gin.H{
				"email":    user.Email,
				"nickname": user.Nickname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireErrorMatch(t, recorder.Body, util.ErrRequired("password"))
			},
		},
		{
			name: "InvalidPasswordType",
			body: gin.H{
				"email":    user.Email,
				"password": 1,
				"nickname": user.Nickname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireErrorMatch(t, recorder.Body, util.ErrType("password", "string"))
			},
		},
		{
			name: "RequiredNickname",
			body: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireErrorMatch(t, recorder.Body, util.ErrRequired("nickname"))
			},
		},
		{
			name: "TooLongNickname",
			body: gin.H{
				"email":    user.Email,
				"password": password,
				"nickname": util.CreateRandomString(51),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireErrorMatch(t, recorder.Body, util.ErrMax("nickname", "50"))
			},
		},
		{
			name: "InvalidNicknameType",
			body: gin.H{
				"email":    user.Email,
				"password": password,
				"nickname": 1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireErrorMatch(t, recorder.Body, util.ErrType("nickname", "string"))
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

func TestLogin(t *testing.T) {
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
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireLoginResponseMatch(t, recorder.Body, user)
			},
		},
		{
			name: "RequiredEmail",
			body: gin.H{
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
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
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
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
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireErrorMatch(t, recorder.Body, util.ErrEmail("email"))
			},
		},
		{
			name: "RequiredPassword",
			body: gin.H{
				"email": user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireErrorMatch(t, recorder.Body, util.ErrRequired("password"))
			},
		},
		{
			name: "InvalidPasswordType",
			body: gin.H{
				"email":    user.Email,
				"password": 1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireErrorMatch(t, recorder.Body, util.ErrType("password", "string"))
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

			url := "/api/auth/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			controller.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestRenewAccessToken(t *testing.T) {
	user, _ := createRandomUser(t)

	testCases := []struct {
		name          string
		body          func(refreshToken string) gin.H
		buildStubs    func(tokenMaker token.Maker, store *mockdb.MockStore) string
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: func(refreshToken string) gin.H {
					return gin.H{
					"refresh_token":    refreshToken,
				}
			},
			buildStubs: func(tokenMaker token.Maker, store *mockdb.MockStore) string {
				token, payload, _ := tokenMaker.CreateToken(user.UserID, time.Minute)

				userAgent := uarand.GetRandom()
				buf := make([]byte, 4)
				ip := rand.Uint32()
				binary.LittleEndian.PutUint32(buf, ip)
				clientIp := net.IP(buf).To4().String()

				session := db.Session{
					SessionID: payload.ID,
					UserID: user.UserID,
					RefreshToken: token,
					UserAgent: userAgent,
					ClientIp: clientIp,
					IsBlocked: false,
					ExpiredAt: time.Now().Add(time.Minute),
				}

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(session, nil)

				return token
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireMatchRenewAccessTokenResponse(t, recorder.Body)
			},
		},
		{
			name: "OK",
			body: func(refreshToken string) gin.H {
					return gin.H{}
			},
			buildStubs: func(tokenMaker token.Maker, store *mockdb.MockStore) string {
				token, _, _ := tokenMaker.CreateToken(user.UserID, time.Minute)

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)

				return token
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireErrorMatch(t, recorder.Body, util.ErrRequired("refresh_token"))
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

			refreshToken := tc.buildStubs(controller.service.TokenMaker, store)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body(refreshToken))
			require.NoError(t, err)

			url := "/api/auth/renew-access-token"
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

func requireLoginResponseMatch(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var rsp dto.LoginResponse
	err = json.Unmarshal(data, &rsp)
	require.NoError(t, err)

	require.NotEmpty(t, rsp.AccessToken)
	require.NotEmpty(t, rsp.RefreshToken)
	require.Equal(t, user.Email, rsp.User.Email)
	require.Equal(t, user.Nickname, rsp.User.Nickname)
	require.WithinDuration(t, user.CreatedAt, rsp.User.CreatedAt, time.Second)
	require.WithinDuration(t, user.UpdatedAt, rsp.User.UpdatedAt, time.Second)
}

func requireMatchRenewAccessTokenResponse(t *testing.T, body *bytes.Buffer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var rsp dto.RenewAccessTokenResponse
	err = json.Unmarshal(data, &rsp)
	require.NoError(t, err)

	require.NotEmpty(t, rsp.AccessToken)
}