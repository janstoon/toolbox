package bricks

import "context"

// Bag is a container which items are non-returnable once pulled. Implementation of Pull must be thread-safe.
type Bag[Value any] interface {
	// Pull takes out an item in case of existence. It should return ErrReachedEnd when nothing remained to return.
	Pull(ctx context.Context) (*Value, error)
}
