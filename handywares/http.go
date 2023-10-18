package handywares

import (
	"log"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/janstoon/toolbox/tricks"
)

type HttpMiddlewareStack middleware.Builder

type (
	pr                            func()
	PanicRecoverHttpMiddlewareOpt = tricks.InPlaceOption[pr]
)

func PrhmoWrapInError(err error) PanicRecoverHttpMiddlewareOpt {
	return func(s *pr) {

	}
}

func (stk *HttpMiddlewareStack) PushPanicRecover(options ...PanicRecoverHttpMiddlewareOpt) *HttpMiddlewareStack {
	//var x pr
	//x = func() {}
	//tricks.ApplyOptions(&x, tricks.Map(func(src PanicRecoverHttpMiddlewareOpt) tricks.Option[pr] {
	//	return src
	//}, options)...)

	stk.Push(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					rw.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(rw, req)
		})
	})

	return stk
}

type BlindLoggerHttpMiddlewareOpt = tricks.InPlaceOption[int]

func (stk *HttpMiddlewareStack) PushBlindLogger(options ...BlindLoggerHttpMiddlewareOpt) *HttpMiddlewareStack {
	stk.Push(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			log.Printf("requested %s", req.URL)

			next.ServeHTTP(rw, req)
		})
	})

	return stk
}

type CorsHttpMiddlewareOpt = tricks.InPlaceOption[int]

func (stk *HttpMiddlewareStack) PushCrossOriginResourceSharingPolicy(
	options ...CorsHttpMiddlewareOpt,
) *HttpMiddlewareStack {
	// headers, methods, origins

	return stk
}

func (stk *HttpMiddlewareStack) Push(mw middleware.Builder) {
	current := *stk
	if current == nil {
		current = middleware.PassthroughBuilder
	}

	*stk = func(next http.Handler) http.Handler {
		return current(mw(next))
	}
}

func (stk *HttpMiddlewareStack) Propagate() {

}
