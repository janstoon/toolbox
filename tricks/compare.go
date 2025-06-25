package tricks

import "cmp"

type Matcher[T any] interface {
	Match(actual T) bool
}

type MatcherFunc[T any] func(actual T) bool

func (f MatcherFunc[T]) Match(actual T) bool {
	return f(actual)
}

func MatchEqual[T comparable](expected T) Matcher[T] {
	return MatcherFunc[T](func(actual T) bool {
		return actual == expected
	})
}

func MatchGreaterThan[T cmp.Ordered](expected T) Matcher[T] {
	return MatcherFunc[T](func(actual T) bool {
		return actual > expected
	})
}

func MatchLesserThan[T cmp.Ordered](expected T) Matcher[T] {
	return MatcherFunc[T](func(actual T) bool {
		return actual < expected
	})
}
