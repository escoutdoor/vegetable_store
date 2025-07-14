package errors

import (
	"errors"
	"fmt"

	"github.com/escoutdoor/vegetable_store/user_service/internal/errors/codes"
)

var (
	ErrJwtTokenExpired       = newError(codes.JwtTokenExpired, "jwt token is already expired")
	ErrInvalidJwtToken       = newError(codes.InvalidJwtToken, "invalid jwt token")
	ErrIncorrectCreadentials = newError(codes.IncorrectCreadentials, "incorrect creadentials")
)

type Error struct {
	Code codes.Code
	Err  error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func newError(code codes.Code, err string) *Error {
	return &Error{
		Code: code,
		Err:  errors.New(err),
	}
}

// Common errors

func UserNotFoundWithID(userID string) *Error {
	msg := fmt.Sprintf("user with id '%s' was not found", userID)
	return newError(codes.UserNotFound, msg)
}

func UserNotFoundWithEmail(email string) *Error {
	msg := fmt.Sprintf("user with email '%s' was not found", email)
	return newError(codes.UserNotFound, msg)
}

func EmailAlreadyExists(email string) *Error {
	msg := fmt.Sprintf("user with email '%s is already exists", email)
	return newError(codes.EmailAlreadyExists, msg)
}
