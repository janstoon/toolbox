package tricks_test

import (
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks"
)

func TestSliceFilter(t *testing.T) {
	assert.Nil(t, tricks.SliceFilter([]int{1, 3, 4, 5}, tricks.MatchEqual(2)))
	assert.Nil(t, tricks.SliceFilter([]int{1, 2, 4, 5, 7, 8}, tricks.MatcherFunc[int](func(src int) bool {
		return src%3 == 0
	})))

	assert.Equal(t, []int{2, 4, 6}, tricks.SliceFilter([]int{1, 2, 3, 4, 5, 6},
		tricks.MatcherFunc[int](func(src int) bool {
			return src%2 == 0
		}),
	))
}

func TestSliceMap(t *testing.T) {
	assert.Nil(t, tricks.SliceMap[[]any, []any](nil, tricks.IdentityMiddleware[any]))
	assert.Nil(t, tricks.SliceMap[[]int, []int]([]int{}, tricks.IdentityMiddleware[int]))

	assert.Equal(t, []int{2, 4, 6, 8, 10, 12}, tricks.SliceMap[[]int, []int]([]int{1, 2, 3, 4, 5, 6},
		func(src int) int {
			return src * 2
		}),
	)

	assert.Equal(t, []string{"value: 1", "value: 2", "value: 3", "value: 4", "value: 5"},
		tricks.SliceMap[[]int, []string]([]int{1, 2, 3, 4, 5}, func(src int) string {
			return fmt.Sprintf("value: %d", src)
		}),
	)
}

func TestSliceReduce(t *testing.T) {
	assert.Equal(t, 720, tricks.SliceReduce([]int{1, 2, 3, 4, 5, 6},
		tricks.ReducerFunc[int, int](func(r int, n int) int {
			if r == 0 {
				return n
			}

			return r * n
		}), 0))

	assert.Equal(t, 7, tricks.SliceReduce([][]int{{1, 2}, {1, 2, 3}, {1, 2}},
		tricks.ReducerFunc[int, []int](func(r int, slice []int) int {
			return r + len(slice)
		}), 0))

	assert.Equal(t, 15, tricks.SliceReduce([]int{1, 2, 3, 4, 5}, tricks.ReduceSum[int](), 0))
}

func TestSliceFind(t *testing.T) {
	v, ok := tricks.SliceFind([]int{1, 2, 3, 4, 5, 6}, tricks.MatchEqual(5))
	assert.True(t, ok)
	assert.Equal(t, 5, v)

	v, ok = tricks.SliceFind([]int{1, 2, 3, 4, 5, 6}, tricks.MatchEqual(50))
	assert.False(t, ok)
	assert.Zero(t, v)

	v, ok = tricks.SliceFind([]int{1, 2, 3, 4, 5, 6}, tricks.MatchGreaterThan(3))
	assert.True(t, ok)
	assert.Equal(t, 4, v)
}

func TestSliceIndex(t *testing.T) {
	assert.Equal(t, 4, tricks.SliceIndex([]int{1, 2, 3, 4, 5, 6}, tricks.MatchEqual(5)))
	assert.Equal(t, 3, tricks.SliceIndex([]int{1, 2, 3, 4, 5, 6}, tricks.MatchGreaterThan(3)))
	assert.Equal(t, tricks.IndexUnavailable, tricks.SliceIndex([]int{1, 2, 3, 4, 5, 6}, tricks.MatchEqual(50)))
}

func TestSliceIndexOf(t *testing.T) {
	assert.Equal(t, 4, tricks.SliceIndexOf(5, []int{1, 2, 3, 4, 5, 6}))
	assert.Equal(t, tricks.IndexUnavailable, tricks.SliceIndexOf(50, []int{1, 2, 3, 4, 5, 6}))
}

func TestSliceFlat(t *testing.T) {
	assert.Nil(t, tricks.SliceFlat[[]any, any]())
	assert.Nil(t, tricks.SliceFlat([]int{}))
	assert.Nil(t, tricks.SliceFlat([]int{}, []int{}))

	assert.Equal(t, []int{1, 2, 1, 3, 4, 2, 3, 4, 1, 7},
		tricks.SliceFlat([]int{1, 2}, []int{1, 3, 4}, []int{2, 3, 4}, []int{1, 7}))

	assert.Equal(t, []int{1, 2, 1, 3, 4, 2, 3, 4, 1, 7},
		tricks.SliceFlat([][]int{{1, 2}, {1, 3, 4}, {2, 3, 4}, {1, 7}}...))
}

func TestSliceClone(t *testing.T) {
	ii := []int{1, 3, 5, 7, 9}
	jj := tricks.SliceClone(ii)
	assert.ElementsMatch(t, ii, jj)

	jj[1] = 2
	assert.NotElementsMatch(t, ii, jj)

	ii[1] = 2
	assert.ElementsMatch(t, ii, jj)

	jj[2] = 4
	assert.NotElementsMatch(t, ii, jj)
	assert.ElementsMatch(t, ii[:2], jj[:2])
	assert.ElementsMatch(t, ii[3:], jj[3:])
}

func TestSliceAll(t *testing.T) {
	assert.True(t, tricks.SliceAll([]int{1, 2, 3, 4, 5}, tricks.MatchLesserThan(10)))

	assert.False(t, tricks.SliceAll([]int{1, 2, 3, 4, 5}, tricks.MatchLesserThan(2)))
}

func TestSliceAny(t *testing.T) {
	assert.True(t, tricks.SliceAny([]int{1, 2, 3, 4, 5}, tricks.MatchLesserThan(10)))
	assert.True(t, tricks.SliceAny([]int{1, 2, 3, 4, 5}, tricks.MatchLesserThan(2)))

	assert.False(t, tricks.SliceAny([]int{1, 2, 3, 4, 5}, tricks.MatchGreaterThan(20)))
}

func TestSliceFlatMap(t *testing.T) {
	assert.Equal(t, []int{0, 0, 1, 0, 1, 4, 0, 1, 4, 9},
		tricks.SliceFlatMap([]int{1, 2, 3, 4}, func(n int) []int {
			dd := make([]int, 0)
			for i := range n {
				dd = append(dd, int(math.Pow(float64(i), 2)))
			}

			return dd
		}),
	)
}

func TestSliceDistinct(t *testing.T) {
	assert.Equal(t, []int{0, 1, 4, 9},
		tricks.SliceDistinct(
			tricks.SliceFlatMap([]int{1, 2, 3, 4}, func(n int) []int {
				dd := make([]int, 0)
				for i := range n {
					dd = append(dd, int(math.Pow(float64(i), 2)))
				}

				return dd
			}),
		),
	)
}

func TestElementAt(t *testing.T) {
	v, ok := tricks.ElementAt(slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9}), 2)
	assert.True(t, ok)
	assert.Equal(t, 5, v)

	v, ok = tricks.ElementAt(slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9}), -10)
	assert.False(t, ok)
	assert.Zero(t, v)

	v, ok = tricks.ElementAt(slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9}), -1)
	assert.False(t, ok)
	assert.Zero(t, v)

	v, ok = tricks.ElementAt(slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9}), 100)
	assert.False(t, ok)
	assert.Zero(t, v)
}

func TestElementAtOrDefault(t *testing.T) {
	assert.Equal(t, 5,
		tricks.ElementAtOrDefault(slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9}), 2, 501))
	assert.Equal(t, 501,
		tricks.ElementAtOrDefault(slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9}), -10, 501))
	assert.Equal(t, 501,
		tricks.ElementAtOrDefault(slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9}), -1, 501))
	assert.Equal(t, 501,
		tricks.ElementAtOrDefault(slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9}), 100, 501))
}

func TestElementAtOrElse(t *testing.T) {
	assert.Equal(t, 5,
		tricks.ElementAtOrElse(slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9}), 2, func() int { return 501 }))
	assert.Equal(t, 501,
		tricks.ElementAtOrElse(slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9}), -10, func() int { return 501 }))
	assert.Equal(t, 501,
		tricks.ElementAtOrElse(slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9}), -1, func() int { return 501 }))
	assert.Equal(t, 501,
		tricks.ElementAtOrElse(slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9}), 100, func() int { return 501 }))
}

func TestSliceDistinctBy(t *testing.T) {
	assert.Equal(t, []int{1, 7, 9, 13, 15}, tricks.SliceDistinctBy([]int{1, 6, 7, 9, 12, 13, 15, 16, 17},
		func(i int) int {
			return i % 5
		}))
}

func TestSkip(t *testing.T) {
	assert.Equal(t, []int{7, 9}, slices.Collect(tricks.Skip(slices.Values([]int{1, 3, 5, 7, 9}), 3)))
}

func TestSkipWhile(t *testing.T) {
	assert.Equal(t, []int{5, 1, 2, 1, 7, 9}, slices.Collect(tricks.SkipWhile(slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9}),
		tricks.MatchLesserThan(5))))
}

func TestTake(t *testing.T) {
	assert.Equal(t, []int{1, 3, 5}, slices.Collect(tricks.Take(slices.Values([]int{1, 3, 5, 7, 9}), 3)))
}

func TestTakeWhile(t *testing.T) {
	assert.Equal(t, []int{1, 3}, slices.Collect(tricks.TakeWhile(slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9}),
		tricks.MatchLesserThan(5))))
}

func TestCount(t *testing.T) {
	assert.Equal(t, 8, tricks.Count(slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9})))
}

func TestSliceLastIndex(t *testing.T) {
	assert.Equal(t, 5, tricks.SliceLastIndex([]int{1, 3, 5, 1, 2, 1, 7, 9}, tricks.MatchEqual(1)))
	assert.Equal(t, 7, tricks.SliceLastIndex([]int{1, 3, 5, 1, 2, 1, 7, 9}, tricks.MatchGreaterThan(1)))
	assert.Equal(t, 5, tricks.SliceLastIndex([]int{1, 3, 5, 1, 2, 1, 7, 9}, tricks.MatchLesserThan(7)))
}

func TestChain(t *testing.T) {
	assert.Equal(t, []int{1, 3, 5, 1, 2, 1, 7, 9, 4, 6, 7, 20, 9}, slices.Collect(tricks.Chain(
		slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9}),
		slices.Values([]int{4, 6, 7, 20, 9}),
	)))
}

func TestExtend(t *testing.T) {
	assert.Equal(t, []int{1, 3, 5, 1, 2, 1, 7, 9, 4, 6, 7, 20, 9}, slices.Collect(tricks.Extend(
		slices.Values([]int{1, 3, 5, 1, 2, 1, 7, 9}),
		[]int{4, 6, 7, 20, 9},
	)))
}
