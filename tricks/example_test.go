package tricks_test

import (
	"fmt"

	"github.com/janstoon/toolbox/tricks"
)

func ExamplePtrVal() {
	fmt.Println(tricks.PtrVal[int](nil))
	fmt.Println(tricks.PtrVal(tricks.ValPtr(4)))

	// Output:
	// 0
	// 4
}

func ExampleCoalesce() {
	fmt.Println(tricks.Coalesce("", "kept", "later"))
	fmt.Println(tricks.Coalesce(0, 0, 3))

	// Output:
	// kept
	// 3
}

func ExampleInsertIfNotExist() {
	hm := map[string]int{}
	tricks.InsertIfNotExist(hm, "a", 1)
	fmt.Println(hm["a"])

	// Output:
	// 1
}

func ExampleMiddlewareStack_Push() {
	var stack tricks.MiddlewareStack[func() string]
	stack = stack.Push(func(next func() string) func() string {
		return func() string { return "a" + next() }
	}).Push(func(next func() string) func() string {
		return func() string { return "b" + next() }
	})

	fmt.Println(stack(func() string { return "c" })())

	// Output:
	// abc
}

func ExampleReduceSum() {
	fmt.Println(tricks.SliceReduce([]int{1, 2, 3}, tricks.ReduceSum[int](), 0))

	// Output:
	// 6
}
