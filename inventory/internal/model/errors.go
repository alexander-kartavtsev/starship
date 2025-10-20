package model

import "github.com/go-faster/errors"

var (
	ErrPartNotFound  error = errors.New("parts not found")
	ErrPartListEmpty error = errors.New("part list is empty")
)
