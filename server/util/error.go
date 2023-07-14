package util

import (
	"fmt"
	"net/http"
)

type DBError struct {
	Name string
	Code string
}

type CustomError struct {
	StatusCode int
	Err        error
}

var (
	DB_UK_ERROR = DBError{Name: "unique_violation", Code: "23505"}
	DB_FK_ERROR = DBError{Name: "foreign_key_violation", Code: "23503"}

	DB_UK_USER_EMAIL    = "users_email_key"
	DB_UK_USER_NICKNAME = "users_nickname_key"

	ErrDuplicateEmail    = CustomError{StatusCode: http.StatusBadRequest, Err: fmt.Errorf("email should be unique")}
	ErrDuplicateNickname = CustomError{StatusCode: http.StatusBadRequest, Err: fmt.Errorf("nickname should be unique")}

	ErrInvalidPassword = CustomError{StatusCode: http.StatusBadRequest, Err: fmt.Errorf("invalid password")}

	ErrNotFoundUser = CustomError{StatusCode: http.StatusNotFound, Err: fmt.Errorf("not found user")}
	ErrNotFoundSession = CustomError{StatusCode: http.StatusNotFound, Err: fmt.Errorf("not found session")}

	ErrBlockedSession = CustomError{StatusCode: http.StatusUnauthorized, Err: fmt.Errorf("blocked session")}
	ErrIncorrectSessionUser = CustomError{StatusCode: http.StatusUnauthorized, Err: fmt.Errorf("incorrect session user")}
	ErrMismatchedSessionToken = CustomError{StatusCode: http.StatusUnauthorized, Err: fmt.Errorf("mismatched session token")}
	ErrExpiredSession = CustomError{StatusCode: http.StatusUnauthorized, Err: fmt.Errorf("expired session")}
)

func ErrType(field string, fieldType string) error {
	return fmt.Errorf("%s should be %s type", field, fieldType)
}

func ErrRequired(field string) error {
	return fmt.Errorf("%s should be required", field)
}

func ErrEmail(field string) error {
	return fmt.Errorf("%s should be email format", field)
}

func ErrMax(field string, param string) error {
	return fmt.Errorf("%s's length should be smaller than or equals to %s", field, param)
}

func NewInternalServerError(err error) CustomError {
	return CustomError{StatusCode: http.StatusInternalServerError, Err: err}
}
