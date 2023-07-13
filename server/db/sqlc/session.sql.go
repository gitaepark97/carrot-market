// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: session.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createSession = `-- name: CreateSession :one
INSERT INTO sessions(
  session_id,
  user_id,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expired_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING sessions.session_id, sessions.user_id, sessions.refresh_token, sessions.user_agent, sessions.client_ip, sessions.is_blocked, sessions.expired_at, sessions.created_at
`

type CreateSessionParams struct {
	SessionID    uuid.UUID `json:"session_id"`
	UserID       int32     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiredAt    time.Time `json:"expired_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, createSession,
		arg.SessionID,
		arg.UserID,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		arg.IsBlocked,
		arg.ExpiredAt,
	)
	var i Session
	err := row.Scan(
		&i.SessionID,
		&i.UserID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiredAt,
		&i.CreatedAt,
	)
	return i, err
}

const getSession = `-- name: GetSession :one
SELECT
  sessions.session_id, sessions.user_id, sessions.refresh_token, sessions.user_agent, sessions.client_ip, sessions.is_blocked, sessions.expired_at, sessions.created_at
FROM sessions
WHERE sessions.session_id = $1
`

func (q *Queries) GetSession(ctx context.Context, sessionID uuid.UUID) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSession, sessionID)
	var i Session
	err := row.Scan(
		&i.SessionID,
		&i.UserID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiredAt,
		&i.CreatedAt,
	)
	return i, err
}
