package tricks_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks"
)

func TestInsertIfNotExist(t *testing.T) {
	hm := make(map[string]int)
	tricks.InsertIfNotExist(hm, "a", 2)
	assert.Equal(t, map[string]int{"a": 2}, hm)
	tricks.InsertIfNotExist(hm, "b", 3)
	assert.Equal(t, map[string]int{"a": 2, "b": 3}, hm)
	require.Panics(t, func() { tricks.InsertIfNotExist(hm, "a", 2) })
	require.Panics(t, func() { tricks.InsertIfNotExist(hm, "a", 20) })
}
