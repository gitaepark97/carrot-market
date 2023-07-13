package token

import "fmt"

var (
	errInvalidToken = fmt.Errorf("token is invalid")
	errExpiredToken = fmt.Errorf("token has expired")
)
