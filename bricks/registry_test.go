package bricks_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/janstoon/toolbox/bricks"
)

func TestRegistry(t *testing.T) {
	ii := bricks.NewRegistry[int]()

	wg := errgroup.Group{}
	for i := 0; i < 5; i++ {
		func(value int) {
			wg.Go(func() error {
				return ii.Register(fmt.Sprintf("I#%d", value+1), value)
			})
		}(i)
	}

	require.NoError(t, wg.Wait())
	assert.Equal(t, map[string]int{"I#1": 0, "I#2": 1, "I#3": 2, "I#4": 3, "I#5": 4}, ii.All())

	v, err := ii.Get("I#5")
	require.NoError(t, err)
	assert.Equal(t, 4, v)

	v, err = ii.Get("I#6")
	require.Error(t, err)
	require.ErrorIs(t, err, bricks.ErrNotFound)
	assert.Zero(t, v)

	err = ii.Register("I#1", 100)
	require.Error(t, err)
	require.ErrorIs(t, err, bricks.ErrAlreadyExists)
	assert.Equal(t, map[string]int{"I#1": 0, "I#2": 1, "I#3": 2, "I#4": 3, "I#5": 4}, ii.All())

	require.NoError(t, ii.Register("I#6", 100))
	assert.Equal(t, map[string]int{"I#1": 0, "I#2": 1, "I#3": 2, "I#4": 3, "I#5": 4, "I#6": 100}, ii.All())

	v, err = ii.Get("I#6")
	require.NoError(t, err)
	assert.Equal(t, 100, v)
}
