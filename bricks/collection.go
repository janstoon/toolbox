package bricks

import (
	"cmp"
	"context"
)

// Collection is the set of objects
type Collection[T any] interface {
	Contains(e T) bool
}

type ClosedInterval[T cmp.Ordered] [2]T // [min, max]

func (tt ClosedInterval[T]) Contains(e T) bool {
	return e >= tt[0] && e <= tt[1]
}

type OpenInterval[T cmp.Ordered] [2]T // (min, max)

func (tt OpenInterval[T]) Contains(e T) bool {
	return e > tt[0] && e < tt[1]
}

type LeftOpenInterval[T cmp.Ordered] [2]T // (min, max]

func (tt LeftOpenInterval[T]) Contains(e T) bool {
	return e > tt[0] && e <= tt[1]
}

type RightOpenInterval[T cmp.Ordered] [2]T // [min, max)

func (tt RightOpenInterval[T]) Contains(e T) bool {
	return e >= tt[0] && e < tt[1]
}

// Bag is a container which items are non-returnable once pulled. Implementation of Pull must be thread-safe.
type Bag[Value any] interface {
	// Pull takes out an item in case of existence. It should return ErrReachedEnd when nothing remained to return.
	Pull(ctx context.Context) (*Value, error)
}
