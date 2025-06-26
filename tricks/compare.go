package tricks

import "cmp"

// Matcher performs matching against value of type T. It's usable in filtering, searching and checking tasks which need
// a predicator to check each individual element.
type Matcher[T any] interface {
	// Match should return true if actual meets the criteria and false if it doesn't.
	Match(actual T) bool
}

// MatcherFunc creates a Matcher from a function which takes value of type T and tells if it meets the criteria or not.
type MatcherFunc[T any] func(actual T) bool

func (f MatcherFunc[T]) Match(actual T) bool {
	return f(actual)
}

// MatchEqual matches actual if it's equal to expected
// actual == expected
func MatchEqual[T comparable](expected T) Matcher[T] {
	return MatcherFunc[T](func(actual T) bool {
		return actual == expected
	})
}

// MatchGreaterThan matches actual if it's greater than expected
// actual > expected
func MatchGreaterThan[T cmp.Ordered](expected T) Matcher[T] {
	return MatcherFunc[T](func(actual T) bool {
		return actual > expected
	})
}

// MatchLesserThan matches actual if it's lesser than expected
// actual < expected
func MatchLesserThan[T cmp.Ordered](expected T) Matcher[T] {
	return MatcherFunc[T](func(actual T) bool {
		return actual < expected
	})
}
