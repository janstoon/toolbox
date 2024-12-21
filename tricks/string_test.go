package tricks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks"
)

func TestStringToRunes(t *testing.T) {
	assert.Equal(t, []rune{'a'}, tricks.StringToRunes("a"))
	assert.Equal(t, []rune{'a', 'b', 'c', 'd'}, tricks.StringToRunes("abcd"))
	assert.Equal(t, []rune{'a', 'b', ' ', 'c', 'd'}, tricks.StringToRunes("ab cd"))
	assert.Equal(t, []rune{'a', 'b', 'b', 'a'}, tricks.StringToRunes("abba"))
}

func TestIsEmptyString(t *testing.T) {
	tcc := []struct {
		s      string
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

	for k, tc := range tcc {
		result := tricks.IsEmptyString(tc.s)
		assert.Equalf(t, tc.result, result, "test case #%d '%s': expected %v, got %v", k, tc.s, tc.result, result)
	}
}
