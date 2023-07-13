package token

import (
	"testing"
	"time"

	"github.com/gitaepark/carrot-market/util"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	tokenMaker, err := NewJWTMaker(util.CreateRandomString(32))
	require.NoError(t, err)

	userID := util.CreateRandomInt32(1, 30)
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := tokenMaker.CreateToken(userID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = tokenMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, userID, payload.UserID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	tokenMaker, err := NewJWTMaker(util.CreateRandomString(32))
	require.NoError(t, err)

	userID := util.CreateRandomInt32(1, 30)
	duration := -time.Minute

	token, payload, err := tokenMaker.CreateToken(userID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = tokenMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, errExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	userID := util.CreateRandomInt32(1, 30)
	duration := time.Minute

	payload, err := NewPayload(userID, duration)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.CreateRandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, errInvalidToken.Error())
	require.Nil(t, payload)
}
