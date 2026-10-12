package bricks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/bricks"
)

func TestQueueAllStopsEarly(t *testing.T) {
	var q bricks.Queue[int]
	q.Enqueue(1, 2, 3)

	var got []int
	for v := range q.All() {
		got = append(got, v)
		if v == 2 {
			break
		}
	}

	assert.Equal(t, []int{1, 2}, got)
}
