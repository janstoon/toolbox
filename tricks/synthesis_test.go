package tricks_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks"
)

type coalesceTestCase[T comparable] struct {
	in  []T
	out T
}

func (tc coalesceTestCase[T]) meet(t *testing.T) {
	assert.Equal(t, tc.out, tricks.Coalesce(tc.in...), "%+v", tc.in)
}

type coalesceTestCases[T comparable] []coalesceTestCase[T]

func (cc coalesceTestCases[T]) meet(t *testing.T) {
	for _, tc := range cc {
		tc.meet(t)
	}
}

func TestCoalesce(t *testing.T) {
	coalesceTestCases[int]{
		{in: []int{}, out: 0},
		{in: []int{0}, out: 0},
		{in: []int{1}, out: 1},
		{in: []int{0, 0}, out: 0},
		{in: []int{1, 0}, out: 1},
		{in: []int{1, 2}, out: 1},
		{in: []int{0, 2}, out: 2},
		{in: []int{0, 0, 0}, out: 0},
		{in: []int{1, 0, 0}, out: 1},
		{in: []int{1, 2, 0}, out: 1},
		{in: []int{1, 2, 3}, out: 1},
		{in: []int{1, 0, 3}, out: 1},
		{in: []int{0, 2, 0}, out: 2},
		{in: []int{0, 2, 3}, out: 2},
		{in: []int{0, 0, 3}, out: 3},
	}.meet(t)

	coalesceTestCases[*int]{
		{in: []*int{}, out: nil},
		{in: []*int{nil}, out: nil},
		{in: []*int{nil, tricks.ValPtr(5)}, out: tricks.ValPtr(5)},
		{in: []*int{nil, nil, tricks.ValPtr(5)}, out: tricks.ValPtr(5)},
	}.meet(t)

	coalesceTestCases[string]{
		{in: []string{}, out: ""},
		{in: []string{""}, out: ""},
		{in: []string{"1st"}, out: "1st"},
		{in: []string{"", ""}, out: ""},
		{in: []string{"1st", ""}, out: "1st"},
		{in: []string{"1st", "2nd"}, out: "1st"},
		{in: []string{"", "2nd"}, out: "2nd"},
		{in: []string{"", "", ""}, out: ""},
		{in: []string{"1st", "", ""}, out: "1st"},
		{in: []string{"1st", "2nd", ""}, out: "1st"},
		{in: []string{"1st", "2nd", "3rd"}, out: "1st"},
		{in: []string{"1st", "", "3rd"}, out: "1st"},
		{in: []string{"", "2nd", ""}, out: "2nd"},
		{in: []string{"", "2nd", "3rd"}, out: "2nd"},
		{in: []string{"", "", "3rd"}, out: "3rd"},
	}.meet(t)

	coalesceTestCases[time.Duration]{
		{in: []time.Duration{}, out: 0 * time.Second},
		{in: []time.Duration{0 * time.Second}, out: 0 * time.Second},
		{in: []time.Duration{1 * time.Second}, out: 1 * time.Second},
		{in: []time.Duration{0 * time.Second, 0 * time.Second}, out: 0 * time.Second},
		{in: []time.Duration{1 * time.Second, 0 * time.Second}, out: 1 * time.Second},
		{in: []time.Duration{1 * time.Second, 2 * time.Second}, out: 1 * time.Second},
		{in: []time.Duration{0 * time.Second, 2 * time.Second}, out: 2 * time.Second},
	}.meet(t)
}
