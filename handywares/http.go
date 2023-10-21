package handywares

import (
	"log"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/janstoon/toolbox/tricks"
	"github.com/rs/cors"
)

type HttpMiddlewareStack middleware.Builder

type PanicRecoverHttpMiddlewareOpt = tricks.InPlaceOption[any]

func (stk *HttpMiddlewareStack) PushPanicRecover(options ...PanicRecoverHttpMiddlewareOpt) *HttpMiddlewareStack {
	return stk.Push(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					rw.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(rw, req)
		})
	})
}

type BlindLoggerHttpMiddlewareOpt = tricks.InPlaceOption[int]

func (stk *HttpMiddlewareStack) PushBlindLogger(options ...BlindLoggerHttpMiddlewareOpt) *HttpMiddlewareStack {
	return stk.Push(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			log.Printf("requested %s", req.URL)

			next.ServeHTTP(rw, req)
		})
	})
}

type CorsHttpMiddlewareOpt = tricks.InPlaceOption[cors.Options]

func (stk *HttpMiddlewareStack) PushCrossOriginResourceSharingPolicy(
	options ...CorsHttpMiddlewareOpt,
) *HttpMiddlewareStack {
	// todo: headers, methods, origins
	cfg := cors.Options{}
	cfg = tricks.PtrVal(tricks.ApplyOptions(&cfg, tricks.Map(func(src CorsHttpMiddlewareOpt) tricks.Option[cors.Options] {
		return src
	}, options)...))

	return stk.Push(cors.New(cfg).Handler)
}

func (stk *HttpMiddlewareStack) Push(mw middleware.Builder) *HttpMiddlewareStack {
	current := *stk
	if current == nil {
		current = middleware.PassthroughBuilder
	}

	*stk = func(next http.Handler) http.Handler {
		return current(mw(next))
	}

	return stk
}

func (stk *HttpMiddlewareStack) NotNil() middleware.Builder {
	if *stk == nil {
		return middleware.PassthroughBuilder
	}

	return middleware.Builder(*stk)
}
