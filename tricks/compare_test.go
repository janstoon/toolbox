package tricks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks"
)

func TestMatchEmptyString(t *testing.T) {
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
		result := tricks.MatchEmptyString().Match(tc.s)
		assert.Equalf(t, tc.result, result, "test case #%d '%s': expected %v, got %v", k, tc.s, tc.result, result)
	}
}
