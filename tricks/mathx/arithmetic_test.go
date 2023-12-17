package mathx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks/mathx"
)

func TestGcd(t *testing.T) {
	assert.Equal(t, 1, mathx.Gcd(1))
	assert.Equal(t, 5, mathx.Gcd(5))
	assert.Equal(t, 6, mathx.Gcd(6))
	assert.Equal(t, 8, mathx.Gcd(8))

	assert.Equal(t, 1, mathx.Gcd(1, 1))
	assert.Equal(t, 5, mathx.Gcd(5, 5))
	assert.Equal(t, 6, mathx.Gcd(6, 6))
	assert.Equal(t, 8, mathx.Gcd(8, 8))
	assert.Equal(t, 2, mathx.Gcd(2, 6))
	assert.Equal(t, 3, mathx.Gcd(3, 9))
	assert.Equal(t, 2, mathx.Gcd(8, 6))
	assert.Equal(t, 2, mathx.Gcd(6, 8))
	assert.Equal(t, 1, mathx.Gcd(2, 3))
	assert.Equal(t, 1, mathx.Gcd(3, 2))
	assert.Equal(t, 1, mathx.Gcd(5, 13))
	assert.Equal(t, 1, mathx.Gcd(11, 13))
	assert.Equal(t, 3, mathx.Gcd(3, 21))
	assert.Equal(t, 14, mathx.Gcd(42, 28))

	assert.Equal(t, 1, mathx.Gcd(1, 1, 1))
	assert.Equal(t, 5, mathx.Gcd(5, 5, 5))
	assert.Equal(t, 6, mathx.Gcd(6, 6, 6))
	assert.Equal(t, 8, mathx.Gcd(8, 8, 8))
	assert.Equal(t, 2, mathx.Gcd(6, 8, 8))
	assert.Equal(t, 2, mathx.Gcd(2, 4, 8))
	assert.Equal(t, 3, mathx.Gcd(3, 9, 21))
	assert.Equal(t, 3, mathx.Gcd(6, 9, 21))
	assert.Equal(t, 7, mathx.Gcd(42, 28, 21))
}

func TestLcm(t *testing.T) {
	assert.Equal(t, 1, mathx.Lcm(1))
	assert.Equal(t, 5, mathx.Lcm(5))
	assert.Equal(t, 6, mathx.Lcm(6))
	assert.Equal(t, 8, mathx.Lcm(8))

	assert.Equal(t, 1, mathx.Lcm(1, 1))
	assert.Equal(t, 5, mathx.Lcm(5, 5))
	assert.Equal(t, 6, mathx.Lcm(6, 6))
	assert.Equal(t, 8, mathx.Lcm(8, 8))
	assert.Equal(t, 6, mathx.Lcm(2, 6))
	assert.Equal(t, 9, mathx.Lcm(3, 9))
	assert.Equal(t, 24, mathx.Lcm(8, 6))
	assert.Equal(t, 24, mathx.Lcm(6, 8))
	assert.Equal(t, 6, mathx.Lcm(2, 3))
	assert.Equal(t, 6, mathx.Lcm(3, 2))
	assert.Equal(t, 65, mathx.Lcm(5, 13))
	assert.Equal(t, 143, mathx.Lcm(11, 13))
	assert.Equal(t, 21, mathx.Lcm(3, 21))
	assert.Equal(t, 84, mathx.Lcm(42, 28))

	assert.Equal(t, 1, mathx.Lcm(1, 1, 1))
	assert.Equal(t, 5, mathx.Lcm(5, 5, 5))
	assert.Equal(t, 6, mathx.Lcm(6, 6, 6))
	assert.Equal(t, 8, mathx.Lcm(8, 8, 8))
	assert.Equal(t, 24, mathx.Lcm(6, 8, 8))
	assert.Equal(t, 8, mathx.Lcm(2, 4, 8))
	assert.Equal(t, 63, mathx.Lcm(3, 9, 21))
	assert.Equal(t, 126, mathx.Lcm(6, 9, 21))
	assert.Equal(t, 84, mathx.Lcm(42, 28, 21))
}
