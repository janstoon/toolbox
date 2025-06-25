package tricks

type Middleware[H any] func(next H) H

// IdentityMiddleware is a transformer which returns input without any modification.
func IdentityMiddleware[H any](h H) H { return h }

type MiddlewareStack[H any] Middleware[H]

func (stk MiddlewareStack[H]) Push(mw Middleware[H]) MiddlewareStack[H] {
	if stk == nil {
		return MiddlewareStack[H](mw)
	}

	return func(next H) H {
		return stk(mw(next))
	}
}
