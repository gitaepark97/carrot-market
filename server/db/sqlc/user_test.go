package db

import (
	"context"
	"testing"
	"time"

	"github.com/gitaepark/carrot-market/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1, _ := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserID, user2.UserID)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Nickname, user2.Nickname)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.UpdatedAt, user2.UpdatedAt, time.Second)
}

func TestGetUserByEmail(t *testing.T) {
	user1, _ := createRandomUser(t)

	user2, err := testQueries.GetUserByEmail(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserID, user2.UserID)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Nickname, user2.Nickname)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.UpdatedAt, user2.UpdatedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	user, _ := createRandomUser(t)

	arg := UpdateUserParams{
		Nickname: util.CreateRandomString(6),
		UserID:   user.UserID,
	}

	updateUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updateUser)

	require.Equal(t, user.UserID, updateUser.UserID)
	require.Equal(t, user.Email, updateUser.Email)
	require.Equal(t, user.HashedPassword, updateUser.HashedPassword)
	require.NotEqual(t, user.Nickname, updateUser.Nickname)
	require.Equal(t, arg.Nickname, updateUser.Nickname)
	require.WithinDuration(t, user.CreatedAt, updateUser.CreatedAt, time.Second)
	require.WithinDuration(t, user.UpdatedAt, updateUser.UpdatedAt, time.Second)
}

func createRandomUser(t *testing.T) (User, string) {
	password := util.CreateRandomString(10)
	hashedPassword, _ := util.HashPassword(password)

	arg := CreateUserParams{
		Email:          util.CreateRandomEmail(),
		HashedPassword: hashedPassword,
		Nickname:       util.CreateRandomString(6),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotZero(t, user.UserID)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Nickname, user.Nickname)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)

	return user, password
}
