package db

import (
	"context"
	"encoding/binary"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/corpix/uarand"
	"github.com/gitaepark/carrot-market/token"
	"github.com/gitaepark/carrot-market/util"
	"github.com/stretchr/testify/require"
)

func TestCreateSession(t *testing.T) {
	createRandomSession(t)
}

func TestGetSession(t *testing.T) {
	session1 := createRandomSession(t)

	session2, err := testQueries.GetSession(context.Background(), session1.SessionID)
	require.NoError(t, err)
	require.NotEmpty(t, session2)

	require.Equal(t, session1.SessionID, session2.SessionID)
	require.Equal(t, session1.UserID, session2.UserID)
	require.Equal(t, session1.RefreshToken, session2.RefreshToken)
	require.Equal(t, session1.UserAgent, session2.UserAgent)
	require.Equal(t, session1.ClientIp, session2.ClientIp)
	require.Equal(t, session1.IsBlocked, session2.IsBlocked)
	require.WithinDuration(t, session1.ExpiredAt, session2.ExpiredAt, time.Second)
	require.WithinDuration(t, session1.CreatedAt, session2.CreatedAt, time.Second)
}

func createRandomSession(t *testing.T) Session {
	maker, _ := token.NewJWTMaker(util.CreateRandomString(32))

	user, _ := createRandomUser(t)

	token, payload, _ := maker.CreateToken(user.UserID, time.Minute)

	buf := make([]byte, 4)
	ip := rand.Uint32()
	binary.LittleEndian.PutUint32(buf, ip)
	clientIp := net.IP(buf).To4().String()

	arg := CreateSessionParams{
		SessionID:    payload.ID,
		UserID:       user.UserID,
		RefreshToken: token,
		UserAgent:    uarand.GetRandom(),
		ClientIp:     clientIp,
		IsBlocked:    false,
		ExpiredAt:    payload.ExpiredAt,
	}

	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, arg.SessionID, session.SessionID)
	require.Equal(t, arg.UserID, session.UserID)
	require.Equal(t, arg.RefreshToken, session.RefreshToken)
	require.Equal(t, arg.UserAgent, session.UserAgent)
	require.Equal(t, arg.ClientIp, session.ClientIp)
	require.Equal(t, arg.IsBlocked, session.IsBlocked)
	require.WithinDuration(t, arg.ExpiredAt, session.ExpiredAt, time.Second)
	require.NotZero(t, session.CreatedAt)

	return session
}
