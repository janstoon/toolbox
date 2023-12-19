package tricks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks"
)

func TestFilter(t *testing.T) {
	assert.Equal(t, []int{2, 4, 6}, tricks.Filter(func(src int) bool {
		return src%2 == 0
	}, []int{1, 2, 3, 4, 5, 6}))
}

func TestMap(t *testing.T) {
	assert.Equal(t, []int{2, 4, 6, 8, 10, 12}, tricks.Map(func(src int) int {
		return src * 2
	}, []int{1, 2, 3, 4, 5, 6}))
}

func TestReduce(t *testing.T) {
	assert.Equal(t, 720, tricks.Reduce(func(r int, n int) int {
		if r == 0 {
			return n
		}

		return r * n
	}, []int{1, 2, 3, 4, 5, 6}))

	assert.Equal(t, 7, tricks.Reduce(func(r int, slice []int) int {
		return r + len(slice)
	}, [][]int{{1, 2}, {1, 2, 3}, {1, 2}}))
}

func TestFind(t *testing.T) {
	assert.Equal(t, tricks.ValPtr(5), tricks.Find(tricks.MatchEqual(5), []int{1, 2, 3, 4, 5, 6}))
	assert.Nil(t, tricks.Find(tricks.MatchEqual(50), []int{1, 2, 3, 4, 5, 6}))
}

func TestFindIndex(t *testing.T) {
	assert.Equal(t, 4, tricks.FindIndex(tricks.MatchEqual(5), []int{1, 2, 3, 4, 5, 6}))
	assert.Equal(t, -1, tricks.FindIndex(tricks.MatchEqual(50), []int{1, 2, 3, 4, 5, 6}))
}

func TestIndexOf(t *testing.T) {
	assert.Equal(t, 4, tricks.IndexOf(5, []int{1, 2, 3, 4, 5, 6}))
	assert.Equal(t, -1, tricks.IndexOf(50, []int{1, 2, 3, 4, 5, 6}))
}

func TestFlat(t *testing.T) {
	assert.Equal(t, []int{1, 2, 1, 3, 4, 2, 3, 4, 1, 7},
		tricks.Flat([]int{1, 2}, []int{1, 3, 4}, []int{2, 3, 4}, []int{1, 7}))
}
