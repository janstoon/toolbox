package bricks

import "errors"

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrInvalidInput   = errors.New("invalid input")
	ErrNotFound       = errors.New("not found")
	ErrInfrastructure = errors.New("infrastructure error")
	ErrReachedEnd     = errors.New("reached end")
)
