package tricks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks"
)

func TestMatchNot(t *testing.T) {
	tcc := []struct {
		v      int
		result bool
	}{
		{1, true},
		{2, true},
		{3, false},
		{4, true},
		{5, true},
		{6, false},
		{7, true},
	}

	matcher := tricks.MatchNot(
		tricks.MatcherFunc[int](func(actual int) bool {
			return actual%3 == 0
		}),
	)
	for k, tc := range tcc {
		result := matcher.Match(tc.v)
		assert.Equalf(t, tc.result, result, "test case #%d '%d': expected %v, got %v", k, tc.v, tc.result, result)
	}
}

func TestMatchAnd(t *testing.T) {
	tcc := []struct {
		v      int
		result bool
	}{
		{1, false},
		{2, false},
		{3, true},
		{4, true},
		{5, true},
		{6, false},
		{7, false},
	}

	matcher := tricks.MatchAnd(
		tricks.MatchGreaterThan(2),
		tricks.MatchLesserThan(6),
	)
	for k, tc := range tcc {
		result := matcher.Match(tc.v)
		assert.Equalf(t, tc.result, result, "test case #%d '%d': expected %v, got %v", k, tc.v, tc.result, result)
	}
}

func TestMatchOr(t *testing.T) {
	tcc := []struct {
		v      int
		result bool
	}{
		{1, true},
		{2, true},
		{3, false},
		{4, false},
		{5, false},
		{6, true},
		{7, true},
	}

	matcher := tricks.MatchOr(
		tricks.MatchLesserThan(3),
		tricks.MatchGreaterThan(5),
	)
	for k, tc := range tcc {
		result := matcher.Match(tc.v)
		assert.Equalf(t, tc.result, result, "test case #%d '%d': expected %v, got %v", k, tc.v, tc.result, result)
	}
}

func TestMatchXor(t *testing.T) {
	tcc := []struct {
		v      int
		result bool
	}{
		{21, true},
		{22, true},
		{23, false},
		{24, true},
		{25, false},
		{26, true},
		{27, true},
		{28, false},
		{29, false},
	}

	matcher := tricks.MatchXor(
		tricks.MatcherFunc[int](func(actual int) bool {
			return actual%2 == 0
		}),
		tricks.MatcherFunc[int](func(actual int) bool {
			return actual%3 == 0
		}),
		tricks.MatcherFunc[int](func(actual int) bool {
			return actual%4 == 0
		}),
	)
	for k, tc := range tcc {
		result := matcher.Match(tc.v)
		assert.Equalf(t, tc.result, result, "test case #%d '%d': expected %v, got %v", k, tc.v, tc.result, result)
	}
}

func TestMatchOneOf(t *testing.T) {
	tcc := []struct {
		v      int
		result bool
	}{
		{21, true},
		{22, true},
		{23, false},
		{24, false},
		{25, false},
		{26, true},
		{27, true},
		{28, false},
		{29, false},
	}

	matcher := tricks.MatchOneOf(
		tricks.MatcherFunc[int](func(actual int) bool {
			return actual%2 == 0
		}),
		tricks.MatcherFunc[int](func(actual int) bool {
			return actual%3 == 0
		}),
		tricks.MatcherFunc[int](func(actual int) bool {
			return actual%4 == 0
		}),
	)
	for k, tc := range tcc {
		result := matcher.Match(tc.v)
		assert.Equalf(t, tc.result, result, "test case #%d '%d': expected %v, got %v", k, tc.v, tc.result, result)
	}
}

func TestMatchEmptyString(t *testing.T) {
	tcc := []struct {
		v      string
		result bool
	}{
		{"", true},
		{" ", true},
		{"   ", true},
		{".", false},
		{" .", false},
		{" . ", false},
		{" _ ", false},
		{" - ", false},
		{"hello", false},
	}

	matcher := tricks.MatchEmptyString()
	for k, tc := range tcc {
		result := matcher.Match(tc.v)
		assert.Equalf(t, tc.result, result, "test case #%d '%s': expected %v, got %v", k, tc.v, tc.result, result)
	}
}
