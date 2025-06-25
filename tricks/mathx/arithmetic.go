package mathx

import (
	"math"
	"math/big"
)

// Gcd calculates the Greatest common divisor of multiple integers
func Gcd(a int, bb ...int) int {
	if len(bb) == 0 {
		return a
	}

	for bb[0] != 0 {
		a, bb[0] = bb[0], a%bb[0]
	}

	return Gcd(a, bb[1:]...)
}

// Lcm calculates the Least common multiple of multiple integers using Gcd
func Lcm(a int, bb ...int) int {
	if len(bb) == 0 {
		return a
	}

	return Lcm(a*bb[0]/Gcd(a, bb[0]), bb[1:]...)
}

// IsPrime checks if the given number (a) is a prime number
func IsPrime(n int) bool {
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}

	return n > 1
}

// PrimeFactors find prime factors of the given number (n)
func PrimeFactors(n int) []int {
	ff := make([]int, 0, 10)

	divisor := 1
	for n > 1 {
		divisor = nextPrime(divisor)
		if n%divisor == 0 {
			ff = append(ff, divisor)

			for n%divisor == 0 {
				n /= divisor
			}
		}
	}

	return ff[:len(ff):len(ff)]
}

// NextPrime returns next prime number greater than the given number (n)
func NextPrime(n int) int {
	if n < 1 {
		n = 1
	}

	return nextPrime(n)
}

// nextPrime returns next prime number greater than the given number (n)
func nextPrime(n int) int {
	for {
		n++

		if IsPrime(n) {
			return n
		}
	}
}

// CoprimesInRange returns list of all natural numbers in closed range [from, to]
// which are coprime to the given number (n)
func CoprimesInRange(n, from, to int) []int {
	cpp := make([]int, 0, 10)
	for i := from; i <= to; i++ {
		if Gcd(n, i) == 1 {
			cpp = append(cpp, i)
		}
	}

	return cpp[:len(cpp):len(cpp)]
}

// MinorCoprimes returns list of all natural numbers which are less than the given number (n) and coprime to it
func MinorCoprimes(n int) []int {
	return CoprimesInRange(n, 1, n-1)
}

// EulerTotient aka φ(n) counts the positive integers up to the given number (n) that are relatively prime to (n)
func EulerTotient(n int) int {
	pff := PrimeFactors(n)

	phi := float64(n)
	for _, pf := range pff {
		phi *= 1 - 1/float64(pf)
	}

	return int(phi)
}

// IsPrimitiveRoot check if (g) is a primitive root of (n)
func IsPrimitiveRoot(n, g int) bool {
	if g < 0 || g >= n || Gcd(g, n) != 1 {
		return false
	}

	bigG, bigN := big.NewInt(int64(g)), big.NewInt(int64(n))

	phi := EulerTotient(n)
	pff := PrimeFactors(phi)
	for _, pf := range pff {
		if big.NewInt(0).Exp(bigG, big.NewInt(int64(phi/pf)), bigN).Int64() == 1 {
			return false
		}
	}

	return true
}

// PrimitiveRoots returns list of all natural numbers which are less than the given number (n) and
// primitive root modulo (n)
func PrimitiveRoots(n int) []int {
	cpp := CoprimesInRange(n, 2, n-1)
	if len(cpp) == 0 {
		return []int{1}
	}

	prr := make([]int, 0, len(cpp))
	for _, cp := range cpp {
		if IsPrimitiveRoot(n, cp) {
			prr = append(prr, cp)
		}
	}

	return prr[:len(prr):len(prr)]
}

// PrimitiveRootsWithRrs returns a map which keys are primitive roots of the given number (n) and values are
// different permutations of least positive reduced residue system modulo (n).
// Each rrs consists of powers of the key, which is a primitive root of (n), modulo (n) in order.
func PrimitiveRootsWithRrs(n int) map[int][]int {
	cpp := CoprimesInRange(n, 2, n-1)
	if len(cpp) == 0 {
		return map[int][]int{1: {1}}
	}

	prr := make(map[int][]int)
	for _, cp := range cpp {
		mm := orderedReducedResidueSystem(n, cp, make([]int, 0, len(cpp)+1))
		if len(mm) == len(cpp)+1 {
			prr[cp] = mm[:len(mm):len(mm)]
		}
	}

	return prr
}

// OrderedReducedResidueSystem returns reduced residue system modulo n
// ordered by powers of g. If g is not a primitive root of n it returns nil
func OrderedReducedResidueSystem(n, g int) []int {
	if !IsPrimitiveRoot(n, g) {
		return nil
	}

	return orderedReducedResidueSystem(n, g, make([]int, 0, EulerTotient(n)))
}

// orderedReducedResidueSystem returns reduced residue system modulo n
// // ordered by powers of g which have to be a primitive root of n
func orderedReducedResidueSystem(n, pr int, rrs []int) []int {
	if rrs == nil {
		rrs = make([]int, 0, 10)
	}

	m := 1
	for i := 1; i < n; i++ {
		m = (m * pr) % n
		rrs = append(rrs, m)

		if m == 1 {
			break
		}
	}

	return rrs[:len(rrs):len(rrs)]
}
