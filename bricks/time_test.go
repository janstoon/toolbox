package bricks_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/bricks"
)

func TestTimespanIncludes(t *testing.T) {
	now := time.Now()

	assert.False(t, bricks.Timespan{Start: now, End: now.Add(5 * time.Hour)}.Contains(now.Add(-2*time.Hour)))
	assert.False(t, bricks.Timespan{Start: now, End: now.Add(5 * time.Hour)}.Contains(now))
	assert.True(t, bricks.Timespan{Start: now, End: now.Add(5 * time.Hour)}.Contains(now.Add(2*time.Hour)))
	assert.False(t, bricks.Timespan{Start: now, End: now.Add(5 * time.Hour)}.Contains(now.Add(6*time.Hour)))
}
