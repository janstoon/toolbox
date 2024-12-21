package tricks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks"
)

func TestCopy(t *testing.T) {
	ii := []int{1, 3, 5, 7, 9}
	jj := tricks.Copy(ii)
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
