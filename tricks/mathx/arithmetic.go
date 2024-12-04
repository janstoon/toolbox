package mathx

import "math"

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
	for {
		if n <= 1 {
			break
		}

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
