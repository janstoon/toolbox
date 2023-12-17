package bricks_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/bricks"
)

func TestClosedIntervalContains(t *testing.T) {
	assert.True(t, bricks.ClosedInterval[int]{0, 3}.Contains(0))
	assert.True(t, bricks.ClosedInterval[int]{0, 3}.Contains(1))
	assert.True(t, bricks.ClosedInterval[int]{0, 3}.Contains(2))
	assert.True(t, bricks.ClosedInterval[int]{0, 3}.Contains(3))

	assert.False(t, bricks.ClosedInterval[int]{1, 3}.Contains(0))
	assert.True(t, bricks.ClosedInterval[int]{1, 3}.Contains(1))
	assert.True(t, bricks.ClosedInterval[int]{1, 3}.Contains(2))
	assert.True(t, bricks.ClosedInterval[int]{1, 3}.Contains(3))
	assert.False(t, bricks.ClosedInterval[int]{1, 3}.Contains(4))

	assert.False(t, bricks.ClosedInterval[float64]{1, 2}.Contains(0))
	assert.True(t, bricks.ClosedInterval[float64]{1, 2}.Contains(1))
	assert.True(t, bricks.ClosedInterval[float64]{1, 2}.Contains(1.5))
	assert.True(t, bricks.ClosedInterval[float64]{1, 2}.Contains(2))
	assert.False(t, bricks.ClosedInterval[float64]{1, 2}.Contains(3))

	assert.True(t, bricks.ClosedInterval[float64]{math.Inf(-1), 2}.Contains(math.Inf(-1)))
	assert.True(t, bricks.ClosedInterval[float64]{math.Inf(-1), 2}.Contains(0))
	assert.True(t, bricks.ClosedInterval[float64]{math.Inf(-1), 2}.Contains(-5))
	assert.True(t, bricks.ClosedInterval[float64]{math.Inf(-1), 2}.Contains(-1000))
	assert.True(t, bricks.LeftOpenInterval[float64]{math.Inf(-1), 2}.Contains(2))
	assert.False(t, bricks.ClosedInterval[float64]{math.Inf(-1), 2}.Contains(5))
	assert.False(t, bricks.ClosedInterval[float64]{math.Inf(-1), 2}.Contains(1000))
	assert.False(t, bricks.ClosedInterval[float64]{math.Inf(-1), 2}.Contains(math.Inf(1)))

	assert.False(t, bricks.ClosedInterval[float64]{1, math.Inf(1)}.Contains(math.Inf(-1)))
	assert.False(t, bricks.ClosedInterval[float64]{1, math.Inf(1)}.Contains(-5))
	assert.False(t, bricks.ClosedInterval[float64]{1, math.Inf(1)}.Contains(0))
	assert.True(t, bricks.ClosedInterval[float64]{1, math.Inf(1)}.Contains(1))
	assert.True(t, bricks.ClosedInterval[float64]{1, math.Inf(1)}.Contains(2))
	assert.True(t, bricks.ClosedInterval[float64]{1, math.Inf(1)}.Contains(5))
	assert.True(t, bricks.ClosedInterval[float64]{1, math.Inf(1)}.Contains(1000))
	assert.True(t, bricks.ClosedInterval[float64]{1, math.Inf(1)}.Contains(math.Inf(1)))
}

func TestOpenIntervalContains(t *testing.T) {
	assert.False(t, bricks.OpenInterval[int]{0, 3}.Contains(0))
	assert.True(t, bricks.OpenInterval[int]{0, 3}.Contains(1))
	assert.True(t, bricks.OpenInterval[int]{0, 3}.Contains(2))
	assert.False(t, bricks.OpenInterval[int]{0, 3}.Contains(3))

	assert.False(t, bricks.OpenInterval[int]{1, 3}.Contains(0))
	assert.False(t, bricks.OpenInterval[int]{1, 3}.Contains(1))
	assert.True(t, bricks.OpenInterval[int]{1, 3}.Contains(2))
	assert.False(t, bricks.OpenInterval[int]{1, 3}.Contains(3))
	assert.False(t, bricks.OpenInterval[int]{1, 3}.Contains(4))

	assert.False(t, bricks.OpenInterval[float64]{1, 2}.Contains(0))
	assert.False(t, bricks.OpenInterval[float64]{1, 2}.Contains(1))
	assert.True(t, bricks.OpenInterval[float64]{1, 2}.Contains(1.5))
	assert.False(t, bricks.OpenInterval[float64]{1, 2}.Contains(2))
	assert.False(t, bricks.OpenInterval[float64]{1, 2}.Contains(3))

	assert.False(t, bricks.OpenInterval[float64]{math.Inf(-1), 2}.Contains(math.Inf(-1)))
	assert.True(t, bricks.OpenInterval[float64]{math.Inf(-1), 2}.Contains(0))
	assert.True(t, bricks.OpenInterval[float64]{math.Inf(-1), 2}.Contains(-5))
	assert.True(t, bricks.OpenInterval[float64]{math.Inf(-1), 2}.Contains(-1000))
	assert.False(t, bricks.OpenInterval[float64]{math.Inf(-1), 2}.Contains(2))
	assert.False(t, bricks.OpenInterval[float64]{math.Inf(-1), 2}.Contains(5))
	assert.False(t, bricks.OpenInterval[float64]{math.Inf(-1), 2}.Contains(1000))
	assert.False(t, bricks.OpenInterval[float64]{math.Inf(-1), 2}.Contains(math.Inf(1)))

	assert.False(t, bricks.OpenInterval[float64]{1, math.Inf(1)}.Contains(math.Inf(-1)))
	assert.False(t, bricks.OpenInterval[float64]{1, math.Inf(1)}.Contains(-5))
	assert.False(t, bricks.OpenInterval[float64]{1, math.Inf(1)}.Contains(0))
	assert.False(t, bricks.OpenInterval[float64]{1, math.Inf(1)}.Contains(1))
	assert.True(t, bricks.OpenInterval[float64]{1, math.Inf(1)}.Contains(2))
	assert.True(t, bricks.OpenInterval[float64]{1, math.Inf(1)}.Contains(5))
	assert.True(t, bricks.OpenInterval[float64]{1, math.Inf(1)}.Contains(1000))
	assert.False(t, bricks.OpenInterval[float64]{1, math.Inf(1)}.Contains(math.Inf(1)))
}

func TestLeftOpenIntervalContains(t *testing.T) {
	assert.False(t, bricks.LeftOpenInterval[int]{0, 3}.Contains(0))
	assert.True(t, bricks.LeftOpenInterval[int]{0, 3}.Contains(1))
	assert.True(t, bricks.LeftOpenInterval[int]{0, 3}.Contains(2))
	assert.True(t, bricks.LeftOpenInterval[int]{0, 3}.Contains(3))

	assert.False(t, bricks.LeftOpenInterval[int]{1, 3}.Contains(0))
	assert.False(t, bricks.LeftOpenInterval[int]{1, 3}.Contains(1))
	assert.True(t, bricks.LeftOpenInterval[int]{1, 3}.Contains(2))
	assert.True(t, bricks.LeftOpenInterval[int]{1, 3}.Contains(3))
	assert.False(t, bricks.LeftOpenInterval[int]{1, 3}.Contains(4))

	assert.False(t, bricks.LeftOpenInterval[float64]{1, 2}.Contains(0))
	assert.False(t, bricks.LeftOpenInterval[float64]{1, 2}.Contains(1))
	assert.True(t, bricks.LeftOpenInterval[float64]{1, 2}.Contains(1.5))
	assert.True(t, bricks.LeftOpenInterval[float64]{1, 2}.Contains(2))
	assert.False(t, bricks.LeftOpenInterval[float64]{1, 2}.Contains(3))

	assert.False(t, bricks.LeftOpenInterval[float64]{math.Inf(-1), 2}.Contains(math.Inf(-1)))
	assert.True(t, bricks.LeftOpenInterval[float64]{math.Inf(-1), 2}.Contains(0))
	assert.True(t, bricks.LeftOpenInterval[float64]{math.Inf(-1), 2}.Contains(-5))
	assert.True(t, bricks.LeftOpenInterval[float64]{math.Inf(-1), 2}.Contains(-1000))
	assert.True(t, bricks.LeftOpenInterval[float64]{math.Inf(-1), 2}.Contains(2))
	assert.False(t, bricks.LeftOpenInterval[float64]{math.Inf(-1), 2}.Contains(5))
	assert.False(t, bricks.LeftOpenInterval[float64]{math.Inf(-1), 2}.Contains(1000))
	assert.False(t, bricks.LeftOpenInterval[float64]{math.Inf(-1), 2}.Contains(math.Inf(1)))

	assert.False(t, bricks.LeftOpenInterval[float64]{1, math.Inf(1)}.Contains(math.Inf(-1)))
	assert.False(t, bricks.LeftOpenInterval[float64]{1, math.Inf(1)}.Contains(-5))
	assert.False(t, bricks.LeftOpenInterval[float64]{1, math.Inf(1)}.Contains(0))
	assert.False(t, bricks.LeftOpenInterval[float64]{1, math.Inf(1)}.Contains(1))
	assert.True(t, bricks.LeftOpenInterval[float64]{1, math.Inf(1)}.Contains(2))
	assert.True(t, bricks.LeftOpenInterval[float64]{1, math.Inf(1)}.Contains(5))
	assert.True(t, bricks.LeftOpenInterval[float64]{1, math.Inf(1)}.Contains(1000))
	assert.True(t, bricks.LeftOpenInterval[float64]{1, math.Inf(1)}.Contains(math.Inf(1)))
}

func TestRightOpenIntervalContains(t *testing.T) {
	assert.True(t, bricks.RightOpenInterval[int]{0, 3}.Contains(0))
	assert.True(t, bricks.RightOpenInterval[int]{0, 3}.Contains(1))
	assert.True(t, bricks.RightOpenInterval[int]{0, 3}.Contains(2))
	assert.False(t, bricks.RightOpenInterval[int]{0, 3}.Contains(3))

	assert.False(t, bricks.RightOpenInterval[int]{1, 3}.Contains(0))
	assert.True(t, bricks.RightOpenInterval[int]{1, 3}.Contains(1))
	assert.True(t, bricks.RightOpenInterval[int]{1, 3}.Contains(2))
	assert.False(t, bricks.RightOpenInterval[int]{1, 3}.Contains(3))
	assert.False(t, bricks.RightOpenInterval[int]{1, 3}.Contains(4))

	assert.False(t, bricks.RightOpenInterval[float64]{1, 2}.Contains(0))
	assert.True(t, bricks.RightOpenInterval[float64]{1, 2}.Contains(1))
	assert.True(t, bricks.RightOpenInterval[float64]{1, 2}.Contains(1.5))
	assert.False(t, bricks.RightOpenInterval[float64]{1, 2}.Contains(2))
	assert.False(t, bricks.RightOpenInterval[float64]{1, 2}.Contains(3))

	assert.True(t, bricks.RightOpenInterval[float64]{math.Inf(-1), 2}.Contains(math.Inf(-1)))
	assert.True(t, bricks.RightOpenInterval[float64]{math.Inf(-1), 2}.Contains(0))
	assert.True(t, bricks.RightOpenInterval[float64]{math.Inf(-1), 2}.Contains(-5))
	assert.True(t, bricks.RightOpenInterval[float64]{math.Inf(-1), 2}.Contains(-1000))
	assert.False(t, bricks.RightOpenInterval[float64]{math.Inf(-1), 2}.Contains(2))
	assert.False(t, bricks.RightOpenInterval[float64]{math.Inf(-1), 2}.Contains(5))
	assert.False(t, bricks.RightOpenInterval[float64]{math.Inf(-1), 2}.Contains(1000))
	assert.False(t, bricks.RightOpenInterval[float64]{math.Inf(-1), 2}.Contains(math.Inf(1)))

	assert.False(t, bricks.RightOpenInterval[float64]{1, math.Inf(1)}.Contains(math.Inf(-1)))
	assert.False(t, bricks.RightOpenInterval[float64]{1, math.Inf(1)}.Contains(-5))
	assert.False(t, bricks.RightOpenInterval[float64]{1, math.Inf(1)}.Contains(0))
	assert.True(t, bricks.RightOpenInterval[float64]{1, math.Inf(1)}.Contains(1))
	assert.True(t, bricks.RightOpenInterval[float64]{1, math.Inf(1)}.Contains(2))
	assert.True(t, bricks.RightOpenInterval[float64]{1, math.Inf(1)}.Contains(5))
	assert.True(t, bricks.RightOpenInterval[float64]{1, math.Inf(1)}.Contains(1000))
	assert.False(t, bricks.RightOpenInterval[float64]{1, math.Inf(1)}.Contains(math.Inf(1)))
}
