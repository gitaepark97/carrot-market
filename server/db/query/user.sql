-- name: CreateUser :one
INSERT INTO users(
  email,
  hashed_password,
  nickname
) VALUES (
  $1, $2, $3
) RETURNING users.*;

-- name: GetUser :one
SELECT
  users.*
FROM users
WHERE users.user_id = $1;

-- name: GetUserByEmail :one
SELECT
  users.*
FROM users
WHERE users.email = $1;

-- name: UpdateUser :one
UPDATE users
SET
  nickname = $1
WHERE users.user_id = $2
RETURNING users.*;