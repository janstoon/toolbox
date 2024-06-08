package tricks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks"
)

func TestMathMiddleware(t *testing.T) {
	type arithmetic func(a, b int) int

	var (
		op0Modifier = func(fn func(a int) int) tricks.Middleware[arithmetic] {
			return func(next arithmetic) arithmetic {
				return func(a, b int) int {
					return next(fn(a), b)
				}
			}
		}

		addOp0By = func(n int) tricks.Middleware[arithmetic] {
			return op0Modifier(func(a int) int {
				return a + n
			})
		}

		subOp0By = func(n int) tricks.Middleware[arithmetic] {
			return op0Modifier(func(a int) int {
				return a - n
			})
		}

		mulOp0By = func(n int) tricks.Middleware[arithmetic] {
			return op0Modifier(func(a int) int {
				return a * n
			})
		}

		divOp0By = func(n int) tricks.Middleware[arithmetic] {
			return op0Modifier(func(a int) int {
				return a / n
			})
		}
	)

	adder := func(a, b int) int { return a + b }

	var mw tricks.MiddlewareStack[arithmetic]

	// 9 + 10
	assert.Equal(t, 19, mw.Push(tricks.IdentityMiddleware[arithmetic])(adder)(9, 10))

	// 5 + 2 - 2 + 10
	assert.Equal(t, 15, mw.Push(addOp0By(2)).Push(subOp0By(2))(adder)(5, 10))

	// 9 * 3 / 2 + 10
	assert.Equal(t, 19, mw.Push(mulOp0By(3)).Push(divOp0By(3))(adder)(9, 10))

	// 5 + 2 * 3 + 10
	assert.Equal(t, 31, mw.Push(addOp0By(2)).Push(mulOp0By(3))(adder)(5, 10))

	// 5 * 3 + 2 + 10
	assert.Equal(t, 27, mw.Push(mulOp0By(3)).Push(addOp0By(2))(adder)(5, 10))
}
