package tricks

import (
	"cmp"
	"strings"
)

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

func MatchTrue[T any]() Matcher[T] {
	return MatcherFunc[T](func(actual T) bool {
		return true
	})
}

func MatchFalse[T any]() Matcher[T] {
	return MatcherFunc[T](func(actual T) bool {
		return false
	})
}

// MatchNot negates a matcher result. It matches what matcher doesn't and vice versa
func MatchNot[T any](matcher Matcher[T]) Matcher[T] {
	return MatcherFunc[T](func(actual T) bool {
		return !matcher.Match(actual)
	})
}

// MatchAnd performs a logic AND on multiple matchers and matches the final result
// It requires all matchers to match
func MatchAnd[T any](matchers ...Matcher[T]) Matcher[T] {
	return MatcherFunc[T](func(actual T) bool {
		for _, matcher := range matchers {
			if !matcher.Match(actual) {
				return false
			}
		}

		return true
	})
}

// MatchOr performs a logic OR on multiple matchers and matches the final result
// It requires at least one matcher to match
func MatchOr[T any](matchers ...Matcher[T]) Matcher[T] {
	return MatcherFunc[T](func(actual T) bool {
		for _, matcher := range matchers {
			if matcher.Match(actual) {
				return true
			}
		}

		return false
	})
}

// MatchXor performs a logic XOR on multiple matchers and matches the final result
// It requires odd number of matchers to match
func MatchXor[T any](matchers ...Matcher[T]) Matcher[T] {
	if len(matchers) == 0 {
		return MatchFalse[T]()
	}

	if len(matchers) == 1 {
		return matchers[0]
	}

	a, b := matchers[0], MatchXor(matchers[1:]...)

	return MatchAnd(
		MatchOr(a, b),
		MatchNot(MatchAnd(a, b)),
	)
}

// MatchOneOf returns a matcher which matches if exactly on of the matchers match
func MatchOneOf[T any](matchers ...Matcher[T]) Matcher[T] {
	return MatcherFunc[T](func(actual T) bool {
		matched := false

		for _, matcher := range matchers {
			if matcher.Match(actual) {
				if matched {
					return false
				}

				matched = true
			}
		}

		return matched
	})
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

// MatchEmptyString matches actual if it's empty string ignoring whitespaces
func MatchEmptyString() Matcher[string] {
	return MatcherFunc[string](func(actual string) bool {
		return len(strings.TrimSpace(actual)) == 0
	})
}
