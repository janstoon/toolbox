package mathx_test

import (
	"fmt"

	"github.com/janstoon/toolbox/tricks/mathx"
)

func ExampleGcd() {
	fmt.Println(mathx.Gcd(12, 18, 30))

	// Output:
	// 6
}

func ExampleLcm() {
	fmt.Println(mathx.Lcm(4, 6))

	// Output:
	// 12
}

func ExampleIsPrime() {
	fmt.Println(mathx.IsPrime(17))
	fmt.Println(mathx.IsPrime(1))

	// Output:
	// true
	// false
}
