package models

import "errors"

type Subject struct {
	ID          int32
	Name        string
	Description string
}

var (
	ErrInvalidDescription = errors.New("invalid description")
)
