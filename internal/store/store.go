package store

import (
	"fmt"
)

type Store interface {
	User() User
	Url() Url
	Alert() Alert
}

type NotFoundError string

func (e NotFoundError) Error() string {
	return string(e)
}

// NewNotFoundError
// Message: {type} with {field}={value} not found
func NewNotFoundError(typ, field string, value any) error {
	return NotFoundError(fmt.Sprintf("%s with %s=%v not found", typ, field, value))
}

func NewNotFoundErrorM(message string) error {
	return NotFoundError(message)
}

type DuplicateError string

func (e DuplicateError) Error() string {
	return string(e)
}

func NewDuplicateError(typ, field string, value any) error {
	return DuplicateError(fmt.Sprintf("%s with %s=%v already exists", typ, field, value))
}
