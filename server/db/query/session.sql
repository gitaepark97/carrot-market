-- name: CreateSession :one
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
) RETURNING sessions.*;

-- name: GetSession :one
SELECT
  sessions.*
FROM sessions
WHERE sessions.session_id = $1;