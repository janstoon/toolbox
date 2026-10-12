package tricks_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks"
)

func TestSequenceGaps(t *testing.T) {
	seq := slices.Values([]int{1, 2, 2, 3})

	i, ok := tricks.IndexOf(slices.All([]int{4, 5, 4}), 4)
	assert.True(t, ok)
	assert.Equal(t, 0, i)
	_, ok = tricks.IndexOf(slices.All([]int{1}), 9)
	assert.False(t, ok)

	i, ok = tricks.LastIndexOf(slices.All([]int{4, 5, 4}), 4)
	assert.True(t, ok)
	assert.Equal(t, 2, i)
	_, ok = tricks.LastIndexOf(slices.All([]int{1}), 9)
	assert.False(t, ok)
	assert.Equal(t, tricks.IndexUnavailable, tricks.SliceLastIndexOf(9, []int{1, 2}))
	assert.Equal(t, 2, tricks.SliceLastIndexOf(1, []int{1, 2, 1}))

	assert.Equal(t, []int{2, 2}, tricks.OrganicSliceFilter([]int{1, 2, 2, 3}, tricks.MatchEqual(2)))
	assert.Nil(t, tricks.OrganicSliceFilter([]int{1}, tricks.MatchFalse[int]()))
	assert.Equal(t, []int{2, 4}, tricks.OrganicSliceMap[[]int, []int]([]int{1, 2}, func(n int) int { return n * 2 }))
	var empty []int
	assert.Nil(t, tricks.OrganicSliceMap[[]int, []int](empty, func(n int) int { return n }))

	assert.Equal(t, []int{2}, tricks.SliceFlatMap([]int{1, 2}, func(n int) []int {
		if n == 1 {
			return nil
		}

		return []int{n}
	}))

	assert.Equal(t, []int{1, 2, 3}, slices.Collect(tricks.Flat(slices.Values([]int{1}), slices.Values([]int{2, 3}))))
	assert.Empty(t, slices.Collect(tricks.Flat[int]()))

	first, ok := tricks.First(seq)
	assert.True(t, ok)
	assert.Equal(t, 1, first)
	_, ok = tricks.First(slices.Values([]int{}))
	assert.False(t, ok)

	last, ok := tricks.Last(slices.Values([]int{7, 8}))
	assert.True(t, ok)
	assert.Equal(t, 8, last)
	_, ok = tricks.Last(slices.Values([]int{}))
	assert.False(t, ok)

	var walked []int
	tricks.ForEach(slices.Values([]int{1, 2}), func(n int) {
		walked = append(walked, n)
	})
	assert.Equal(t, []int{1, 2}, walked)
}
