package errors

import (
	"fmt"
	"net/http"
)

type ApplicationError interface {
	GetHTTPStatus() int
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

type ItemAlreadyExistsError struct {
	ItemId uint64
}

func (err *ItemAlreadyExistsError) Error() string {
	return fmt.Sprintf("Item with id=%d already exists", err.ItemId)
}

func (err *ItemAlreadyExistsError) GetHTTPStatus() int {
	return http.StatusBadRequest
}
