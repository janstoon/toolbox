package bricks

import (
	"context"
)

// Compensator undoes work that already happened when a later step fails.
type Compensator interface {
	// Compensate undoes the failed step described by err.
	Compensate(ctx context.Context, err error) error
}

// CompensatorFunc adapts a function to Compensator.
type CompensatorFunc func(ctx context.Context, err error) error

// Compensate calls cf.
func (cf CompensatorFunc) Compensate(ctx context.Context, err error) error {
	return cf(ctx, err)
}

type errorAndCompensator struct {
	error
	Compensator
}

// CompensatorAsError joins err with compensator so errors.As can recover the compensator.
func CompensatorAsError(compensator Compensator, err error) error {
	return errorAndCompensator{
		error:       err,
		Compensator: compensator,
	}
}

// Unwrap returns the wrapped error.
func (eac errorAndCompensator) Unwrap() error {
	return eac.error
}
