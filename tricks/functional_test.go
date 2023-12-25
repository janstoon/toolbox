package tricks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks"
)

func TestFilter(t *testing.T) {
	assert.Nil(t, tricks.Filter([]int{1, 3, 4, 5}, tricks.MatchEqual(2)))
	assert.Nil(t, tricks.Filter([]int{1, 2, 4, 5, 7, 8}, func(src int) bool {
		return src%3 == 0
	}))

	assert.Equal(t, []int{2, 4, 6}, tricks.Filter([]int{1, 2, 3, 4, 5, 6}, func(src int) bool {
		return src%2 == 0
	}))
}

func TestMap(t *testing.T) {
	assert.Nil(t, tricks.Map(nil, func(src any) any {
		return src
	}))
	assert.Equal(t, []int{}, tricks.Map([]int{}, func(src int) int {
		return src
	}))

	assert.Equal(t, []int{2, 4, 6, 8, 10, 12}, tricks.Map([]int{1, 2, 3, 4, 5, 6}, func(src int) int {
		return src * 2
	}))
}

func TestReduce(t *testing.T) {
	assert.Equal(t, 720, tricks.Reduce([]int{1, 2, 3, 4, 5, 6}, func(r int, n int) int {
		if r == 0 {
			return n
		}

		return r * n
	}))

	assert.Equal(t, 7, tricks.Reduce([][]int{{1, 2}, {1, 2, 3}, {1, 2}}, func(r int, slice []int) int {
		return r + len(slice)
	}))
}

func TestFind(t *testing.T) {
	assert.Equal(t, tricks.ValPtr(5), tricks.Find([]int{1, 2, 3, 4, 5, 6}, tricks.MatchEqual(5)))
	assert.Nil(t, tricks.Find([]int{1, 2, 3, 4, 5, 6}, tricks.MatchEqual(50)))
}

func TestFindIndex(t *testing.T) {
	assert.Equal(t, 4, tricks.FindIndex([]int{1, 2, 3, 4, 5, 6}, tricks.MatchEqual(5)))
	assert.Equal(t, -1, tricks.FindIndex([]int{1, 2, 3, 4, 5, 6}, tricks.MatchEqual(50)))
}

func TestIndexOf(t *testing.T) {
	assert.Equal(t, 4, tricks.IndexOf(5, []int{1, 2, 3, 4, 5, 6}))
	assert.Equal(t, -1, tricks.IndexOf(50, []int{1, 2, 3, 4, 5, 6}))
}

func TestFlat(t *testing.T) {
	assert.Nil(t, tricks.Flat[any]())
	assert.Nil(t, tricks.Flat([]int{}))
	assert.Nil(t, tricks.Flat([]int{}, []int{}))

	assert.Equal(t, []int{1, 2, 1, 3, 4, 2, 3, 4, 1, 7},
		tricks.Flat([]int{1, 2}, []int{1, 3, 4}, []int{2, 3, 4}, []int{1, 7}))

	assert.Equal(t, []int{1, 2, 1, 3, 4, 2, 3, 4, 1, 7},
		tricks.Flat([][]int{{1, 2}, {1, 3, 4}, {2, 3, 4}, {1, 7}}...))
}
