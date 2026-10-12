package bricks

import (
	"cmp"
	"context"
)

// Collection is the set of objects
type Collection[T any] interface {
	Contains(e T) bool
}

// ClosedInterval is the inclusive range [min, max]. Index 0 is the lower bound and index 1 is the upper bound.
type ClosedInterval[T cmp.Ordered] [2]T

// Contains reports whether e is inside the closed interval.
func (tt ClosedInterval[T]) Contains(e T) bool {
	return e >= tt[0] && e <= tt[1]
}

// OpenInterval is the exclusive range (min, max). Index 0 is the lower bound and index 1 is the upper bound.
type OpenInterval[T cmp.Ordered] [2]T

// Contains reports whether e is inside the open interval.
func (tt OpenInterval[T]) Contains(e T) bool {
	return e > tt[0] && e < tt[1]
}

// LeftOpenInterval is the range (min, max]. Index 0 is the lower bound and index 1 is the upper bound.
type LeftOpenInterval[T cmp.Ordered] [2]T

// Contains reports whether e is inside the left-open interval.
func (tt LeftOpenInterval[T]) Contains(e T) bool {
	return e > tt[0] && e <= tt[1]
}

// RightOpenInterval is the range [min, max). Index 0 is the lower bound and index 1 is the upper bound.
type RightOpenInterval[T cmp.Ordered] [2]T

// Contains reports whether e is inside the right-open interval.
func (tt RightOpenInterval[T]) Contains(e T) bool {
	return e >= tt[0] && e < tt[1]
}

// Bag is a container which items are non-returnable once pulled. Implementation of Pull must be thread-safe.
type Bag[Value any] interface {
	// Pull takes out an item in case of existence. It should return ErrReachedEnd when nothing remained to return.
	Pull(ctx context.Context) (*Value, error)
}
