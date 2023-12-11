package base62_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks/wire/base62"
)

func TestEncode(t *testing.T) {
	assert.Equal(t, []byte("K6zrUM5k"), base62.DefaultEndec.Encode([]byte("Pouyan")))
	assert.Equal(t, []byte("JG"), base62.DefaultEndec.Encode([]byte("M")))
	assert.Equal(t, []byte("JM4"), base62.DefaultEndec.Encode([]byte("Ma")))
	assert.Equal(t, []byte("JM5k"), base62.DefaultEndec.Encode([]byte("Man")))
	assert.Equal(t, []byte("JM5kUG"), base62.DefaultEndec.Encode([]byte("Many")))
}
