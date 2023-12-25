package tricks_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks"
)

func TestPtrVal(t *testing.T) {
	assert.Zero(t, tricks.PtrVal[int](nil))
	assert.Zero(t, tricks.PtrVal[string](nil))
	assert.Zero(t, tricks.PtrVal[bool](nil))
	assert.Zero(t, tricks.PtrVal[[]int](nil))
	assert.Zero(t, tricks.PtrVal[map[string]int](nil))

	assert.Equal(t, 0, tricks.PtrVal[int](nil))
	assert.Equal(t, "", tricks.PtrVal[string](nil))
	assert.False(t, tricks.PtrVal[bool](nil))
	assert.Nil(t, tricks.PtrVal[[]int](nil))
	assert.Nil(t, tricks.PtrVal[map[string]int](nil))

	assert.Equal(t, 0, tricks.PtrVal(tricks.ValPtr(0)))
	assert.Equal(t, 3, tricks.PtrVal(tricks.ValPtr(3)))
	assert.Equal(t, []int{2, 4, 5}, tricks.PtrVal(tricks.ValPtr([]int{2, 4, 5})))

	assert.Equal(t, []int{2, 4, 5}, tricks.PtrVal(&[]int{2, 4, 5}))
}

func TestValPtr(t *testing.T) {
	str := "yalda"
	assert.Exactly(t, &str, tricks.ValPtr(str))
	assert.NotSame(t, &str, tricks.ValPtr(str))
	assert.NotSame(t, tricks.ValPtr(str), tricks.ValPtr(str))

	var iarr []int
	assert.Nil(t, iarr)
	assert.Exactly(t, &iarr, tricks.ValPtr(iarr))
}

func TestPtrPtr(t *testing.T) {
	assert.Nil(t, tricks.PtrPtr(nil, strconv.Itoa))

	str := "neda"
	assert.Exactly(t, tricks.ValPtr(4), tricks.PtrPtr(&str, func(src string) int {
		return len(src)
	}))
}
