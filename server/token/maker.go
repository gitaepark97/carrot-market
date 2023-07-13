package token

import "time"

type Maker interface {
	CreateToken(userID int32, duration time.Duration) (string, *Payload, error)
	VerifyToken(tokenString string) (*Payload, error)
}
