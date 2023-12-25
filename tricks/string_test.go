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
