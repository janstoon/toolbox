package tricks

// Middleware is simply a function that generates H which ought to be a handler.
// It's a wrapper that takes the next H as input and decides to call it or not in the generated H anywhere desired.
// H can be http.Handler, grpc.UnaryServerInterceptor, queue message handler or any other task handler.
// Simplest Middleware is Identity which directly returns the next H as the generated H aka passthrough middleware.
type Middleware[H any] func(next H) H

// MiddlewareStack is a mechanism to chain multiple Middleware(s) and output a single entrypoint.
// Middleware(s) should all accept and generate same handler type.
// Think of it as http server middleware builder with Middleware[http.Handler].
// MiddlewareStack is originally a function and calling it with root handler as input gives a handler which runs all
// pushed Middleware(s) one-by-one in FIFO order.
type MiddlewareStack[H any] Middleware[H]

// Push appends the Middleware to the stack and returns the new MiddlewareStack.
// It's a function composer and Middleware(s) gets called in order of Push that means calling further Push
// on returning MiddlewareStack appends the Middleware to ehe end,
// and it gets called after all previously pushed Middleware(s).
func (stk MiddlewareStack[H]) Push(mw Middleware[H]) MiddlewareStack[H] {
	if stk == nil {
		return MiddlewareStack[H](mw)
	}

	return func(next H) H {
		return stk(mw(next))
	}
}
