package bricks

import "errors"

var (
	ErrReachedEnd = errors.New("reached end")

	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrNotFound           = errors.New("not found")

	ErrInfrastructure = errors.New("infrastructure error")
	ErrNotImplemented = errors.New("not implemented")
)
