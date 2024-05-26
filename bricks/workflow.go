package bricks

import (
	"context"
)

type Compensator interface {
	Compensate(ctx context.Context) error
}

type CompensatorFunc func(ctx context.Context) error

func (cf CompensatorFunc) Compensate(ctx context.Context) error {
	return cf(ctx)
}

type errorAndCompensator struct {
	error
	Compensator
}

func CompensatorAsError(compensator Compensator, err error) error {
	return errorAndCompensator{
		error:       err,
		Compensator: compensator,
	}
}

func (eac errorAndCompensator) Unwrap() error {
	return eac.error
}
