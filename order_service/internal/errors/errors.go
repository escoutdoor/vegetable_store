package errors

import (
	"fmt"
	"strings"

	"github.com/escoutdoor/vegetable_store/order_service/internal/errors/codes"
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

func OrderNotFound(orderID string) *Error {
	return newError(codes.OrderNotFound, fmt.Errorf("order with id '%s' was not found", orderID))
}

func InsufficientInventory(vegetableID string, requested, available float32) *Error {
	return newError(codes.InsufficientInventory, fmt.Errorf("insufficient inventory for vegetable %s: requested %.2f kg, available %.2f kg", vegetableID, requested, available))
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
