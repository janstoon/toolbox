package tricks

import "cmp"

// Reducer folds a sequence of elements into one accumulator.
type Reducer[A, E any] interface {
	// Reduce combines accumulator with element and returns the next accumulator.
	Reduce(accumulator A, element E) A
}

// ReducerFunc adapts a function to Reducer.
type ReducerFunc[A, E any] func(accumulator A, element E) A

// Reduce calls f.
func (f ReducerFunc[A, E]) Reduce(accumulator A, element E) A {
	return f(accumulator, element)
}

// ReduceSum returns a reducer that adds each element to the accumulator.
func ReduceSum[T cmp.Ordered]() Reducer[T, T] {
	return ReducerFunc[T, T](func(accumulator, element T) T {
		return accumulator + element
	})
}

// Coalesce returns left-most non-zero value
// It's like cmp.Or
func Coalesce[T comparable](tt ...T) T {
	var zero T
	if len(tt) == 0 {
		return zero
	}

	if tt[0] != zero || len(tt) == 1 {
		return tt[0]
	}

	return Coalesce[T](tt[1:]...)
}
