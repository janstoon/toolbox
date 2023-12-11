package mathx

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
