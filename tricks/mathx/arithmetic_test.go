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
	assert.Equal(t, 1, mathx.Gcd(999999000001, 2))

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
	assert.Equal(t, 1999998000002, mathx.Lcm(999999000001, 2))

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

func TestIsPrime(t *testing.T) {
	assert.True(t, mathx.IsPrime(2))
	assert.True(t, mathx.IsPrime(3))
	assert.True(t, mathx.IsPrime(5))
	assert.True(t, mathx.IsPrime(7))
	assert.True(t, mathx.IsPrime(11))
	assert.True(t, mathx.IsPrime(13))
	assert.True(t, mathx.IsPrime(97))
	assert.True(t, mathx.IsPrime(127))
	assert.True(t, mathx.IsPrime(281))
	assert.True(t, mathx.IsPrime(389))
	assert.True(t, mathx.IsPrime(499))
	assert.True(t, mathx.IsPrime(8191))
	assert.True(t, mathx.IsPrime(524287))
	assert.True(t, mathx.IsPrime(6700417))
	assert.True(t, mathx.IsPrime(999999000001))

	assert.False(t, mathx.IsPrime(1))
	assert.False(t, mathx.IsPrime(4))
	assert.False(t, mathx.IsPrime(6))
	assert.False(t, mathx.IsPrime(8))
	assert.False(t, mathx.IsPrime(9))
	assert.False(t, mathx.IsPrime(10))
	assert.False(t, mathx.IsPrime(12))
	assert.False(t, mathx.IsPrime(68))
	assert.False(t, mathx.IsPrime(72))
	assert.False(t, mathx.IsPrime(90))
	assert.False(t, mathx.IsPrime(93))
	assert.False(t, mathx.IsPrime(100))
	assert.False(t, mathx.IsPrime(121))
	assert.False(t, mathx.IsPrime(6700415))
	assert.False(t, mathx.IsPrime(6700419))
	assert.False(t, mathx.IsPrime(999999000000))
	assert.False(t, mathx.IsPrime(999999000002))
	assert.False(t, mathx.IsPrime(999999000004))

	assert.False(t, mathx.IsPrime(0))
	assert.False(t, mathx.IsPrime(-1))
	assert.False(t, mathx.IsPrime(-2))
	assert.False(t, mathx.IsPrime(-3))
	assert.False(t, mathx.IsPrime(-4))
	assert.False(t, mathx.IsPrime(-5))
	assert.False(t, mathx.IsPrime(-6))
	assert.False(t, mathx.IsPrime(-7))
}

func TestPrimeFactors(t *testing.T) {
	assert.Equal(t, []int{2, 3}, mathx.PrimeFactors(6))
	assert.Equal(t, []int{2, 5}, mathx.PrimeFactors(10))
	assert.Equal(t, []int{2, 3, 5}, mathx.PrimeFactors(30))
	assert.Equal(t, []int{3, 7}, mathx.PrimeFactors(21))
	assert.Equal(t, []int{2, 3, 5, 13}, mathx.PrimeFactors(390))
	assert.Equal(t, []int{2, 3, 5, 7, 13}, mathx.PrimeFactors(2730))
	assert.Equal(t, []int{2, 3, 5, 7, 13}, mathx.PrimeFactors(40950))
	assert.Equal(t, []int{2, 3, 5, 7, 13, 17}, mathx.PrimeFactors(46410))

	assert.Equal(t, []int{2}, mathx.PrimeFactors(4))
	assert.Equal(t, []int{2}, mathx.PrimeFactors(8))
	assert.Equal(t, []int{2}, mathx.PrimeFactors(16))
	assert.Equal(t, []int{2}, mathx.PrimeFactors(32))
	assert.Equal(t, []int{2}, mathx.PrimeFactors(64))
	assert.Equal(t, []int{2}, mathx.PrimeFactors(128))
	assert.Equal(t, []int{2}, mathx.PrimeFactors(256))
	assert.Equal(t, []int{2}, mathx.PrimeFactors(512))
	assert.Equal(t, []int{2}, mathx.PrimeFactors(1024))
	assert.Equal(t, []int{2}, mathx.PrimeFactors(1<<15))

	assert.Equal(t, []int{}, mathx.PrimeFactors(1))
	assert.Equal(t, []int{2}, mathx.PrimeFactors(2))
	assert.Equal(t, []int{3}, mathx.PrimeFactors(3))
	assert.Equal(t, []int{5}, mathx.PrimeFactors(5))
	assert.Equal(t, []int{7}, mathx.PrimeFactors(7))
	assert.Equal(t, []int{11}, mathx.PrimeFactors(11))
	assert.Equal(t, []int{13}, mathx.PrimeFactors(13))

	assert.Equal(t, []int{}, mathx.PrimeFactors(0))
	assert.Equal(t, []int{}, mathx.PrimeFactors(-1))
	assert.Equal(t, []int{}, mathx.PrimeFactors(-2))
	assert.Equal(t, []int{}, mathx.PrimeFactors(-3))
	assert.Equal(t, []int{}, mathx.PrimeFactors(-4))
	assert.Equal(t, []int{}, mathx.PrimeFactors(-5))
	assert.Equal(t, []int{}, mathx.PrimeFactors(-6))
	assert.Equal(t, []int{}, mathx.PrimeFactors(-7))
	assert.Equal(t, []int{}, mathx.PrimeFactors(-100))
	assert.Equal(t, []int{}, mathx.PrimeFactors(-121))
	assert.Equal(t, []int{}, mathx.PrimeFactors(-127))
}

func TestNextPrime(t *testing.T) {
	assert.Equal(t, 2, mathx.NextPrime(0))
	assert.Equal(t, 2, mathx.NextPrime(1))
	assert.Equal(t, 3, mathx.NextPrime(2))
	assert.Equal(t, 5, mathx.NextPrime(3))
	assert.Equal(t, 5, mathx.NextPrime(4))
	assert.Equal(t, 7, mathx.NextPrime(5))
	assert.Equal(t, 999999000001, mathx.NextPrime(999999000000))

	assert.Equal(t, 2, mathx.NextPrime(-1))
	assert.Equal(t, 2, mathx.NextPrime(-2))
	assert.Equal(t, 2, mathx.NextPrime(-3))
	assert.Equal(t, 2, mathx.NextPrime(-4))
	assert.Equal(t, 2, mathx.NextPrime(-5))
	assert.Equal(t, 2, mathx.NextPrime(-100))
}

func TestCoprimesInRange(t *testing.T) {
	assert.Equal(t, []int{3, 7}, mathx.CoprimesInRange(10, 2, 7))
	assert.Equal(t, []int{7, 9, 11, 13, 17, 19}, mathx.CoprimesInRange(10, 5, 20))
}

func TestMinorCoprimes(t *testing.T) {
	assert.Equal(t, []int{1, 5}, mathx.MinorCoprimes(6))
	assert.Equal(t, []int{1, 3, 7, 9}, mathx.MinorCoprimes(10))
}

func TestEulerTotient(t *testing.T) {
	assert.Equal(t, 1, mathx.EulerTotient(1))
	assert.Equal(t, 1, mathx.EulerTotient(2))
	assert.Equal(t, 2, mathx.EulerTotient(3))
	assert.Equal(t, 2, mathx.EulerTotient(4))
	assert.Equal(t, 4, mathx.EulerTotient(5))
	assert.Equal(t, 6, mathx.EulerTotient(7))
	assert.Equal(t, 4, mathx.EulerTotient(10))
	assert.Equal(t, 12, mathx.EulerTotient(13))
}

func TestIsPrimitiveRoot(t *testing.T) {
	assert.True(t, mathx.IsPrimitiveRoot(2, 1))
	assert.True(t, mathx.IsPrimitiveRoot(3, 2))
	assert.True(t, mathx.IsPrimitiveRoot(7, 3))
	assert.True(t, mathx.IsPrimitiveRoot(7, 5))
	assert.True(t, mathx.IsPrimitiveRoot(10, 3))
	assert.True(t, mathx.IsPrimitiveRoot(10, 7))
	assert.True(t, mathx.IsPrimitiveRoot(50, 3))
	assert.True(t, mathx.IsPrimitiveRoot(50, 13))
	assert.True(t, mathx.IsPrimitiveRoot(50, 17))
	assert.True(t, mathx.IsPrimitiveRoot(50, 23))
	assert.True(t, mathx.IsPrimitiveRoot(50, 27))
	assert.True(t, mathx.IsPrimitiveRoot(50, 33))
	assert.True(t, mathx.IsPrimitiveRoot(50, 37))
	assert.True(t, mathx.IsPrimitiveRoot(50, 47))

	assert.False(t, mathx.IsPrimitiveRoot(2, 0))
	assert.False(t, mathx.IsPrimitiveRoot(2, 2))
	assert.False(t, mathx.IsPrimitiveRoot(2, 3))
	assert.False(t, mathx.IsPrimitiveRoot(2, 5))
	assert.False(t, mathx.IsPrimitiveRoot(3, 1))
	assert.False(t, mathx.IsPrimitiveRoot(3, 3))
	assert.False(t, mathx.IsPrimitiveRoot(3, 4))
	assert.False(t, mathx.IsPrimitiveRoot(3, 5))
	assert.False(t, mathx.IsPrimitiveRoot(7, 1))
	assert.False(t, mathx.IsPrimitiveRoot(7, 2))
	assert.False(t, mathx.IsPrimitiveRoot(7, 4))
	assert.False(t, mathx.IsPrimitiveRoot(7, 6))
	assert.False(t, mathx.IsPrimitiveRoot(7, 7))
	assert.False(t, mathx.IsPrimitiveRoot(7, 8))
	assert.False(t, mathx.IsPrimitiveRoot(7, 9))
	assert.False(t, mathx.IsPrimitiveRoot(7, 10))
	assert.False(t, mathx.IsPrimitiveRoot(50, 41))
	assert.False(t, mathx.IsPrimitiveRoot(50, 42))
	assert.False(t, mathx.IsPrimitiveRoot(50, 43))
}

func TestPrimitiveRoots(t *testing.T) {
	assert.Equal(t, []int{1}, mathx.PrimitiveRoots(2))
	assert.Equal(t, []int{2}, mathx.PrimitiveRoots(3))
	assert.Equal(t, []int{3, 5}, mathx.PrimitiveRoots(7))
	assert.Equal(t, []int{3, 7}, mathx.PrimitiveRoots(10))
	assert.Equal(t, []int{3, 5, 6, 7, 10, 11, 12, 14}, mathx.PrimitiveRoots(17))

	nn := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 20, 21, 22, 23, 24, 25, 49, 50, 95, 96, 97, 98, 99, 100,
	}
	for _, n := range nn {
		prr := mathx.PrimitiveRoots(n)
		if len(prr) > 0 {
			assert.Lenf(t, prr, mathx.EulerTotient(mathx.EulerTotient(n)), "number of primitive roots (%d) mismatch", n)
		}
	}
}

func TestPrimitiveRootsWithRrs(t *testing.T) {
	assert.Equal(t, map[int][]int{
		1: {1},
	}, mathx.PrimitiveRootsWithRrs(2))

	assert.Equal(t, map[int][]int{
		2: {2, 1},
	}, mathx.PrimitiveRootsWithRrs(3))

	assert.Equal(t, map[int][]int{
		2: {2, 4, 3, 1},
		3: {3, 4, 2, 1},
	}, mathx.PrimitiveRootsWithRrs(5))

	assert.Equal(t, map[int][]int{
		3: {3, 2, 6, 4, 5, 1},
		5: {5, 4, 6, 2, 3, 1},
	}, mathx.PrimitiveRootsWithRrs(7))

	assert.Equal(t, map[int][]int{
		3: {3, 9, 7, 1},
		7: {7, 9, 3, 1},
	}, mathx.PrimitiveRootsWithRrs(10))

	assert.Equal(t, map[int][]int{
		2: {2, 4, 8, 5, 10, 9, 7, 3, 6, 1},
		6: {6, 3, 7, 9, 10, 5, 8, 4, 2, 1},
		7: {7, 5, 2, 3, 10, 4, 6, 9, 8, 1},
		8: {8, 9, 6, 4, 10, 3, 2, 5, 7, 1},
	}, mathx.PrimitiveRootsWithRrs(11))

	assert.Equal(t, map[int][]int{
		3:  {3, 9, 10, 13, 5, 15, 11, 16, 14, 8, 7, 4, 12, 2, 6, 1},
		5:  {5, 8, 6, 13, 14, 2, 10, 16, 12, 9, 11, 4, 3, 15, 7, 1},
		6:  {6, 2, 12, 4, 7, 8, 14, 16, 11, 15, 5, 13, 10, 9, 3, 1},
		7:  {7, 15, 3, 4, 11, 9, 12, 16, 10, 2, 14, 13, 6, 8, 5, 1},
		10: {10, 15, 14, 4, 6, 9, 5, 16, 7, 2, 3, 13, 11, 8, 12, 1},
		11: {11, 2, 5, 4, 10, 8, 3, 16, 6, 15, 12, 13, 7, 9, 14, 1},
		12: {12, 8, 11, 13, 3, 2, 7, 16, 5, 9, 6, 4, 14, 15, 10, 1},
		14: {14, 9, 7, 13, 12, 15, 6, 16, 3, 8, 10, 4, 5, 2, 11, 1},
	}, mathx.PrimitiveRootsWithRrs(17))

	nn := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 20, 21, 22, 23, 24, 25, 49, 50, 95, 96, 97, 98, 99, 100,
	}
	for _, n := range nn {
		prr, phi := mathx.PrimitiveRootsWithRrs(n), mathx.EulerTotient(n)
		for _, pra := range prr {
			assert.Len(t, pra, phi)
			for _, prb := range prr {
				assert.ElementsMatch(t, pra, prb)
			}
		}
	}
}

func TestOrderedReducedResidueSystem(t *testing.T) {
	assert.Equal(t, []int{2, 4, 3, 1}, mathx.OrderedReducedResidueSystem(5, 2))
	assert.Equal(t, []int{3, 4, 2, 1}, mathx.OrderedReducedResidueSystem(5, 3))
	assert.Equal(t, []int{
		63, 10, 95, 100, 94, 37, 84, 49, 91, 62, 54, 85, 5, 101, 50, 47, 72, 42, 78, 99, 31, 27, 96, 56, 104, 25, 77, 36,
		21, 39, 103, 69, 67, 48, 28, 52, 66, 92, 18, 64, 73, 105, 88, 87, 24, 14, 26, 33, 46, 9, 32, 90, 106, 44, 97, 12,
		7, 13, 70, 23, 58, 16, 45, 53, 22, 102, 6, 57, 60, 35, 65, 29, 8, 76, 80, 11, 51, 3, 82, 30, 71, 86, 68, 4, 38, 40,
		59, 79, 55, 41, 15, 89, 43, 34, 2, 19, 20, 83, 93, 81, 74, 61, 98, 75, 17, 1,
	}, mathx.OrderedReducedResidueSystem(107, 63))
}
