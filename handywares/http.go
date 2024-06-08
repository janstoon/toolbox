package handywares

import (
	"errors"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/tricks"
	"github.com/rs/cors"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

type HttpMiddlewareStack middleware.Builder

type PanicRecoverHttpMiddlewareOpt = tricks.InPlaceOption[any]

func (stk *HttpMiddlewareStack) PushPanicRecover(options ...PanicRecoverHttpMiddlewareOpt) *HttpMiddlewareStack {
	return stk.Push(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("paniced %+v\n", r)
					debug.PrintStack()

					rw.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(rw, req)
		})
	})
}

type BlindLoggerHttpMiddlewareOpt = tricks.InPlaceOption[any]

func (stk *HttpMiddlewareStack) PushBlindLogger(
	mctx *middleware.Context, options ...BlindLoggerHttpMiddlewareOpt,
) *HttpMiddlewareStack {
	return stk.Push(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			log.Printf(
				"HTTP|%s/%s %s [%s] %s %s %s\n",
				req.RemoteAddr,
				req.Referer(),
				req.UserAgent(),
				req.Method,
				req.Host,
				req.URL,
				req.Proto,
			)

			next.ServeHTTP(rw, req)
		})
	})
}

type CorsHttpMiddlewareOpt = tricks.InPlaceOption[cors.Options]

func (stk *HttpMiddlewareStack) PushCrossOriginResourceSharingPolicy(
	options ...CorsHttpMiddlewareOpt,
) *HttpMiddlewareStack {
	cfg := cors.Options{}
	cfg = tricks.PtrVal(tricks.ApplyOptions(&cfg,
		tricks.Map(options, func(src CorsHttpMiddlewareOpt) tricks.Option[cors.Options] {
			return src
		})...))

	return stk.Push(cors.New(cfg).Handler)
}

var CorsAllowOrigins = func(origins ...string) CorsHttpMiddlewareOpt {
	return func(s *cors.Options) {
		s.AllowedOrigins = origins
	}
}

var CorsAllowMethods = func(methods ...string) CorsHttpMiddlewareOpt {
	return func(s *cors.Options) {
		s.AllowedMethods = methods
	}
}

var CorsAllowHeaders = func(headers ...string) CorsHttpMiddlewareOpt {
	return func(s *cors.Options) {
		s.AllowedHeaders = headers
	}
}

var CorsDebug = func(debug bool) CorsHttpMiddlewareOpt {
	return func(s *cors.Options) {
		s.Debug = debug
	}
}

type HttpRouteTester func(route *middleware.MatchedRoute) bool

func PassthroughHttpRouteTester(success bool) HttpRouteTester {
	return func(_ *middleware.MatchedRoute) bool {
		return success
	}
}

func CombineHttpRouteTesters(tt ...HttpRouteTester) HttpRouteTester {
	switch len(tt) {
	case 1:
		return tt[0]

	case 0:
		return PassthroughHttpRouteTester(true)
	}

	return func(route *middleware.MatchedRoute) bool {
		return tt[0](route) && CombineHttpRouteTesters(tt[1:]...)(route)
	}
}

type OtelHmw struct {
	tracer trace.Tracer
	mctx   *middleware.Context

	namePrefix  string
	routeTester HttpRouteTester
}

type OpenTelemetryHttpMiddlewareOpt = tricks.Option[OtelHmw]

func OtelHttpSpanNamePrefix(prefix string) OpenTelemetryHttpMiddlewareOpt {
	return tricks.OutOfPlaceOption[OtelHmw](func(hmw OtelHmw) OtelHmw {
		hmw.namePrefix = prefix

		return hmw
	})
}

func OtelHttpRouteTester(tester HttpRouteTester) OpenTelemetryHttpMiddlewareOpt {
	return tricks.OutOfPlaceOption[OtelHmw](func(hmw OtelHmw) OtelHmw {
		hmw.routeTester = tester

		return hmw
	})
}

func OtelHttpOperationIdException(oids ...string) OpenTelemetryHttpMiddlewareOpt {
	return OtelHttpRouteTester(func(route *middleware.MatchedRoute) bool {
		return tricks.IndexOf(route.Operation.ID, oids) < 0
	})
}

func (stk *HttpMiddlewareStack) PushOpenTelemetry(
	tracer trace.Tracer, mctx *middleware.Context, options ...OpenTelemetryHttpMiddlewareOpt,
) *HttpMiddlewareStack {
	hmw := &OtelHmw{
		tracer: tracer,
		mctx:   mctx,
	}
	hmw = tricks.ApplyOptions(hmw, options...)

	return stk.Push(hmw.builder)
}

func (hmw OtelHmw) builder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		if route, match := hmw.match(req); match {
			var span trace.Span
			ctx, span = hmw.tracer.Start(req.Context(), hmw.spanName(route.Operation.ID))
			defer span.End()

			span.SetAttributes(
				semconv.HTTPRequestMethodKey.String(req.Method),
			)
		}

		next.ServeHTTP(rw, req.WithContext(ctx))
	})
}

func (hmw OtelHmw) match(req *http.Request) (*middleware.MatchedRoute, bool) {
	if route, matched := hmw.mctx.LookupRoute(req); matched && hmw.routeTester(route) {
		return route, true
	}

	return nil, false
}

func (hmw OtelHmw) spanName(opId string) string {
	sb := strings.Builder{}
	if len(strings.TrimSpace(hmw.namePrefix)) > 0 {
		sb.WriteString(hmw.namePrefix)
		sb.WriteRune('/')
	}

	sb.WriteString(opId)

	return sb.String()
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

type HttpTripperwareStack = tricks.MiddlewareStack[http.RoundTripper]

func HttpErrorMapperTripperware(
	mapper func(*http.Response, error) (*http.Response, error),
) tricks.Middleware[http.RoundTripper] {
	if mapper == nil {
		panic(errors.Join(bricks.ErrInvalidArgument, errors.New("empty error mapper")))
	}

	return func(next http.RoundTripper) http.RoundTripper {
		return HttpRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return mapper(next.RoundTrip(req))
		})
	}
}

type HttpRoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn HttpRoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
