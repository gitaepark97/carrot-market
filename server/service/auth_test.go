package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	mockdb "github.com/gitaepark/carrot-market/db/mock"
	db "github.com/gitaepark/carrot-market/db/sqlc"
	"github.com/gitaepark/carrot-market/dto"
	"github.com/gitaepark/carrot-market/util"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	user, password := createRandomUser(t)

	testCases := []struct {
		name          string
		reqBody       dto.RegisterRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(rsp dto.UserResponse, err util.CustomError)
	}{
		{
			name: "OK",
			reqBody: dto.RegisterRequest{
				Email:    user.Email,
				Password: password,
				Nickname: user.Nickname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(rsp dto.UserResponse, err util.CustomError) {
				requireMatchUserResponse(t, rsp, user)
				require.Empty(t, err)
			},
		},
		{
			name: "InternalServerError",
			reqBody: dto.RegisterRequest{
				Email:    user.Email,
				Password: password,
				Nickname: user.Nickname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(rsp dto.UserResponse, err util.CustomError) {
				require.Empty(t, rsp)
				requireErrorMatch(t, err, util.NewInternalServerError(sql.ErrConnDone))
			},
		},
		{
			name: "DuplicateEmail",
			reqBody: dto.RegisterRequest{
				Email:    user.Email,
				Password: password,
				Nickname: user.Nickname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: pq.ErrorCode(util.DB_UK_ERROR.Code), Constraint: util.DB_UK_USER_EMAIL})
			},
			checkResponse: func(rsp dto.UserResponse, err util.CustomError) {
				require.Empty(t, rsp)
				requireErrorMatch(t, err, util.ErrDuplicateEmail)
			},
		},
		{
			name: "DuplicateNickname",
			reqBody: dto.RegisterRequest{
				Email:    user.Email,
				Password: password,
				Nickname: user.Nickname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: pq.ErrorCode(util.DB_UK_ERROR.Code), Constraint: util.DB_UK_USER_NICKNAME})
			},
			checkResponse: func(rsp dto.UserResponse, err util.CustomError) {
				require.Empty(t, rsp)
				requireErrorMatch(t, err, util.ErrDuplicateNickname)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			service := newTestService(t, store)

			tc.buildStubs(store)

			rsp, err := service.Register(context.Background(), tc.reqBody)
			tc.checkResponse(rsp, err)
		})
	}
}

func TestLogin(t *testing.T) {
	user, password := createRandomUser(t)

	testCases := []struct {
		name          string
		reqBody       dto.LoginRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(rsp dto.LoginResponse, err util.CustomError)
	}{
		{
			name: "OK",
			reqBody: dto.LoginRequest{
				Email:    user.Email,
				Password: password,
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
			checkResponse: func(rsp dto.LoginResponse, err util.CustomError) {
				requireMatchLoginResponse(t, rsp, user)
				require.Empty(t, err)
			},
		},
		{
			name: "InternalServerError",
			reqBody: dto.LoginRequest{
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rsp dto.LoginResponse, err util.CustomError) {
				require.Empty(t, rsp)
				requireErrorMatch(t, err, util.NewInternalServerError(sql.ErrConnDone))
			},
		},
		{
			name: "NotFoundUser",
			reqBody: dto.LoginRequest{
				Email:    util.CreateRandomEmail(),
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rsp dto.LoginResponse, err util.CustomError) {
				require.Empty(t, rsp)
				requireErrorMatch(t, err, util.ErrNotFoundUser)
			},
		},
		{
			name: "InvalidPassword",
			reqBody: dto.LoginRequest{
				Email:    user.Email,
				Password: util.CreateRandomPassword(),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rsp dto.LoginResponse, err util.CustomError) {
				require.Empty(t, rsp)
				requireErrorMatch(t, err, util.ErrInvalidPassword)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			service := newTestService(t, store)

			tc.buildStubs(store)

			rsp, err := service.Login(context.Background(), tc.reqBody)
			tc.checkResponse(rsp, err)
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

func requireMatchUserResponse(t *testing.T, rsp dto.UserResponse, user db.User) {
	require.Equal(t, rsp.Email, user.Email)
	require.Equal(t, rsp.Nickname, user.Nickname)
	require.WithinDuration(t, rsp.CreatedAt, user.CreatedAt, time.Second)
	require.WithinDuration(t, rsp.UpdatedAt, user.UpdatedAt, time.Second)
}

func requireMatchLoginResponse(t *testing.T, rsp dto.LoginResponse, user db.User) {
	require.NotEmpty(t, rsp.AccessToken)
	require.NotEmpty(t, rsp.RefreshToken)
	require.Equal(t, rsp.User.Email, user.Email)
	require.Equal(t, rsp.User.Nickname, user.Nickname)
	require.WithinDuration(t, rsp.User.CreatedAt, user.CreatedAt, time.Second)
	require.WithinDuration(t, rsp.User.UpdatedAt, user.UpdatedAt, time.Second)
}
