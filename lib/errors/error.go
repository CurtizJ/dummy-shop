package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode uint

const (
	UNKNOWN_ERROR ErrorCode = iota
	ITEM_NOT_FOUND
	ITEM_ALREADY_EXISTS
)

type ApplicationError interface {
	GetHTTPStatus() int
	Code() ErrorCode
}

type ItemNotFoundError struct {
	ItemId uint64
}

func (err *ItemNotFoundError) Error() string {
	return fmt.Sprintf("Item with id=%d not found", err.ItemId)
}

func (err *ItemNotFoundError) GetHTTPStatus() int {
	return http.StatusNotFound
}

func (err *ItemNotFoundError) Code() ErrorCode {
	return ITEM_NOT_FOUND
}

type ItemAlreadyExistsError struct {
	ItemId uint64
}

func (err *ItemAlreadyExistsError) Error() string {
	return fmt.Sprintf("Item with id=%d already exists", err.ItemId)
}

func (err *ItemAlreadyExistsError) GetHTTPStatus() int {
	return http.StatusBadRequest
}

func (err *ItemAlreadyExistsError) Code() ErrorCode {
	return ITEM_ALREADY_EXISTS
}
