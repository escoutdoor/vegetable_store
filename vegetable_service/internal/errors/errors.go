package errors

import (
	"fmt"
	"strings"

	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/errors/codes"
)

type Error struct {
	Code codes.Code
	Err  error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func newError(code codes.Code, err error) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}

// Common errors

func VegetableNotFound(vegetableID string) *Error {
	return newError(codes.VegetableNotFound, fmt.Errorf("vegetable with id '%s' was not found", vegetableID))
}

func VegetablesNotFound(vegetableIDs []string) *Error {
	var msg error
	if len(vegetableIDs) == 1 {
		msg = fmt.Errorf("vegetable with id '%s' was not found", vegetableIDs[0])
	} else {
		msg = fmt.Errorf("vegetables with ids '%s' were not found", strings.Join(vegetableIDs, "', '"))
	}

	return newError(codes.VegetablesNotFound, msg)
}
