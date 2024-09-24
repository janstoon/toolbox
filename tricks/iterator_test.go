package tricks_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks"
)

func TestIterator_All(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}
	it := tricks.NewSliceIterator(intSlice)

	it_f, ok := tricks.NewSliceIterator([]int{1, 2, 3, 4, 5}).Next()
	assert.True(t, ok)
	assert.Equal(t, 1, it_f)

	iterator := it.Iter()
	iterator_cloned := iterator.Clone()
	iterator.All(func(i int) bool {
		return i < 10
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, iterator_cloned.Collect())
	assert.Equal(t, []int{}, iterator.Collect())

	assert.True(t, it.Iter().All(func(i int) bool {
		return i < 10
	}))

	assert.False(t, it.Iter().All(func(i int) bool {
		return i < 5
	}))

	assert.True(t, it.Iter().Any(func(i int) bool {
		return i > 3
	}))

	it_bigger_than_2 := it.Iter().Filter(func(i int) bool {
		return i > 2
	}).Collect()
	assert.Equal(t, []int{3, 4, 5}, it_bigger_than_2)

	it_sum := it.Iter().Reduce(func(acc int, i int) int {
		return acc + i
	})
	assert.Equal(t, 15, it_sum)

	it_string_bullets := tricks.IteratorCast[string](it.Iter().Map(func(i int) any {
		return fmt.Sprintf("value: %d", i)
	})).Collect()
	assert.Equal(t, []string{"value: 1", "value: 2", "value: 3", "value: 4", "value: 5"}, it_string_bullets)

	assert.True(t, it.Iter().Contains(3))
	assert.False(t, it.Iter().Contains(10))
	assert.True(t, it.Iter().ContainsBy(func(i int) bool {
		return i == 3
	}))
	assert.False(t, it.Iter().ContainsBy(func(i int) bool {
		return i == 10
	}))

	v, found := it.Iter().Find(func(i int) bool {
		return i%3 == 0
	})
	assert.Equal(t, 3, v)
	assert.True(t, found)

	v, found = it.Iter().Find(func(i int) bool {
		return i == 10
	})
	assert.Equal(t, 0, v)
	assert.False(t, found)

	indx := it.Iter().FindIndex(func(i int) bool {
		return i%3 == 0
	})
	assert.Equal(t, 2, indx)

	indx = it.Iter().FindIndex(func(i int) bool {
		return i == 10
	})
	assert.Equal(t, -1, indx)

	assert.Equal(t, 1, it.Iter().FlatMap(func(i int) []int {
		return []int{i, i}
	}).LastIndexOf(1))

	flat_map := it.Iter().FlatMap(func(i int) []int {
		return []int{i, i}
	}).Collect()
	assert.Equal(t, []int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5}, flat_map)

	distinct := it.Iter().Distinct().Collect()
	assert.Equal(t, []int{1, 2, 3, 4, 5}, distinct)
	distinct2 := it.Iter().FlatMap(func(i int) []int {
		return []int{i, i}
	}).Distinct().Collect()
	assert.Equal(t, []int{1, 2, 3, 4, 5}, distinct2)
	distinct3 := it.Iter().DistinctBy(func(i int) int {
		return i % 2
	}).Collect()
	assert.Equal(t, []int{1, 2}, distinct3)

	first, has := it.Iter().First()
	assert.True(t, has)
	assert.Equal(t, 1, first)

	last, has2 := it.Iter().Last()
	assert.True(t, has2)
	assert.Equal(t, 5, last)

	it3 := tricks.NewSliceIterator([]int{})
	first3, has3 := it3.Iter().First()
	assert.False(t, has3)
	assert.Equal(t, 0, first3)
	last3, has4 := it3.Iter().Last()
	assert.False(t, has4)
	assert.Equal(t, 0, last3)

	elem_at, haselem := it.Iter().ElementAt(2)
	assert.True(t, haselem)
	assert.Equal(t, 3, elem_at)

	assert.Equal(t, 5, it.Iter().ElementAtOrDefault(10, 5))
	assert.Equal(t, 5, it.Iter().ElementAtOrElse(10, func() int {
		return 5
	}))

	cnt := it.Iter().Count()
	assert.Equal(t, 5, cnt)

	var multiplied_slice []int
	it.Iter().ForEach(func(i int) {
		multiplied_slice = append(multiplied_slice, i*2)
	})
	assert.Equal(t, []int{2, 4, 6, 8, 10}, multiplied_slice)

	assert.Equal(t, -1, it.Iter().IndexOf(10))
	assert.Equal(t, 4, tricks.IteratorCast[int](it.Iter().Map(func(i int) any { return 1 })).LastIndexOf(1))

	skip := it.Iter().Skip(2).Collect()
	assert.Equal(t, []int{3, 4, 5}, skip)

	skip_while := it.Iter().SkipWhile(func(i int) bool {
		return i < 3
	}).Collect()
	assert.Equal(t, []int{3, 4, 5}, skip_while)

	take_2 := it.Iter().Take(2).Collect()
	assert.Equal(t, []int{1, 2}, take_2)

	take_while := it.Iter().TakeWhile(func(i int) bool {
		return i < 3
	}).Collect()
	assert.Equal(t, []int{1, 2}, take_while)

	chained := it.Iter().Chain(tricks.NewSliceIterator([]int{6, 7, 8, 9, 10})).Collect()
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, chained)

	extended := it.Iter().Extend([]int{6, 7, 8, 9, 10}).Collect()
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, extended)

	it2 := tricks.IteratorCast[string](it.Iter().
		Map(func(i int) any {
			return "value: " + strconv.Itoa(i)
		}).
		Map(func(i any) any {
			str, _ := i.(string)

			return str + "!"
		})).
		Collect()

	assert.Equal(t, []string{"value: 1!", "value: 2!", "value: 3!", "value: 4!", "value: 5!"}, it2)
}
