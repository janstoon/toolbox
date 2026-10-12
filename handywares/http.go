package handywares

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"runtime/debug"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/tricks"
	"github.com/rs/cors"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// HttpMiddlewareStack chains HTTP server middleware. The first Push runs first.
type HttpMiddlewareStack = tricks.MiddlewareStack[http.Handler]

// PanicRecoverHttpMiddlewareOpt configures HttpPanicRecoverMiddleware. No options are defined yet.
type PanicRecoverHttpMiddlewareOpt = tricks.MutableOption[any]

// HttpPanicRecoverMiddleware recovers a panic in the next handler, records it on the current span,
// and writes HTTP 500.
func HttpPanicRecoverMiddleware(options ...PanicRecoverHttpMiddlewareOpt) tricks.Middleware[http.Handler] {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					span := trace.SpanFromContext(req.Context())
					span.AddEvent("panic recovered", trace.WithAttributes(
						oaPanicValue.String(fmt.Sprintf("%+v", r)),
						oaDebugStack.String(string(debug.Stack())),
					))

					// todo: translate to proper http status
					rw.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(rw, req)
		})
	}
}

// BlindLoggerHttpMiddlewareOpt configures HttpBlindLoggerMiddleware. No options are defined yet.
type BlindLoggerHttpMiddlewareOpt = tricks.MutableOption[any]

// HttpBlindLoggerMiddleware logs the request line before calling the next handler.
// mctx is accepted for symmetry with other HTTP middleware and is not read.
func HttpBlindLoggerMiddleware(
	mctx *middleware.Context, options ...BlindLoggerHttpMiddlewareOpt,
) tricks.Middleware[http.Handler] {
	return func(next http.Handler) http.Handler {
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
	}
}

// CorsHttpMiddlewareOpt changes the CORS policy used by HttpCrossOriginResourceSharingPolicyMiddleware.
type CorsHttpMiddlewareOpt = tricks.MutableOption[cors.Options]

// HttpCrossOriginResourceSharingPolicyMiddleware applies a CORS policy built from options.
func HttpCrossOriginResourceSharingPolicyMiddleware(options ...CorsHttpMiddlewareOpt) tricks.Middleware[http.Handler] {
	cfg := cors.Options{}
	cfg = tricks.PtrVal(tricks.ApplyOptions(&cfg,
		tricks.SliceMap[[]CorsHttpMiddlewareOpt, []tricks.Option[cors.Options]](options,
			func(src CorsHttpMiddlewareOpt) tricks.Option[cors.Options] {
				return src
			},
		)...))

	return cors.New(cfg).Handler
}

// CorsAllowOrigins sets the origins the CORS policy accepts.
var CorsAllowOrigins = func(origins ...string) CorsHttpMiddlewareOpt {
	return func(s *cors.Options) {
		s.AllowedOrigins = origins
	}
}

// CorsAllowMethods sets the methods the CORS policy accepts.
var CorsAllowMethods = func(methods ...string) CorsHttpMiddlewareOpt {
	return func(s *cors.Options) {
		s.AllowedMethods = methods
	}
}

// CorsAllowHeaders sets the headers the CORS policy accepts.
var CorsAllowHeaders = func(headers ...string) CorsHttpMiddlewareOpt {
	return func(s *cors.Options) {
		s.AllowedHeaders = headers
	}
}

// CorsDebug turns CORS debug logging on or off.
var CorsDebug = func(debug bool) CorsHttpMiddlewareOpt {
	return func(s *cors.Options) {
		s.Debug = debug
	}
}

// HttpRouteTester reports whether a matched OpenAPI route should be handled specially.
type HttpRouteTester func(route *middleware.MatchedRoute) bool

// PassthroughHttpRouteTester returns a tester that ignores the route and always returns success.
func PassthroughHttpRouteTester(success bool) HttpRouteTester {
	return func(_ *middleware.MatchedRoute) bool {
		return success
	}
}

// CombineHttpRouteTesters returns a tester that succeeds only when every tester succeeds.
// No testers means every route succeeds. One tester is returned unchanged.
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

// OtelHmw is the OpenTelemetry HTTP server middleware built by HttpOpenTelemetryMiddleware.
type OtelHmw struct {
	tracer trace.Tracer
	mctx   *middleware.Context

	namePrefix  string
	routeTester HttpRouteTester
}

// OpenTelemetryHttpMiddlewareOpt configures HttpOpenTelemetryMiddleware.
type OpenTelemetryHttpMiddlewareOpt = tricks.Option[OtelHmw]

// OtelHttpSpanNamePrefix sets the text written before the operation id in the span name.
func OtelHttpSpanNamePrefix(prefix string) OpenTelemetryHttpMiddlewareOpt {
	return tricks.ImmutableOption[OtelHmw](func(hmw OtelHmw) OtelHmw {
		hmw.namePrefix = prefix

		return hmw
	})
}

// OtelHttpRouteTester selects which matched routes receive a real span.
// Routes the tester rejects are traced with a no-op tracer.
func OtelHttpRouteTester(tester HttpRouteTester) OpenTelemetryHttpMiddlewareOpt {
	return tricks.ImmutableOption[OtelHmw](func(hmw OtelHmw) OtelHmw {
		hmw.routeTester = tester

		return hmw
	})
}

// OtelHttpOperationIdException traces every matched route except the named operation ids.
func OtelHttpOperationIdException(oids ...string) OpenTelemetryHttpMiddlewareOpt {
	return OtelHttpRouteTester(func(route *middleware.MatchedRoute) bool {
		return tricks.SliceIndexOf(route.Operation.ID, oids) < 0
	})
}

// HttpOpenTelemetryMiddleware starts a server span for requests whose route passes the tester.
// Unmatched requests, and matched requests the tester rejects, use a no-op tracer.
func HttpOpenTelemetryMiddleware(
	tracer trace.Tracer, mctx *middleware.Context, options ...OpenTelemetryHttpMiddlewareOpt,
) tricks.Middleware[http.Handler] {
	hmw := &OtelHmw{
		tracer: tracer,
		mctx:   mctx,
	}
	hmw = tricks.ApplyOptions(hmw, options...)

	return hmw.builder
}

func (hmw OtelHmw) builder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		tracer := noop.NewTracerProvider().Tracer("")

		route, matched := hmw.match(req)
		if matched {
			tracer = hmw.tracer
		}

		oid := "?"
		attrs := []attribute.KeyValue{
			semconv.HTTPRequestMethodKey.String(req.Method),
			semconv.URLScheme(req.URL.Scheme),
			semconv.URLFull(req.URL.String()),
			semconv.NetworkProtocolVersion(req.Proto),
			semconv.UserAgentOriginal(req.UserAgent()),
			semconv.ServerAddress(req.Host),
			semconv.ClientAddress(req.RemoteAddr),
			attribute.Key("http.request.header.referer").String(req.Referer()),
		}

		if route != nil {
			oid = route.Operation.ID
			attrs = append(attrs, semconv.HTTPRoute(route.PathPattern))
		}

		ctx, span := tracer.Start(req.Context(), hmw.spanName(oid), trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		span.SetAttributes(attrs...)

		next.ServeHTTP(rw, req.WithContext(ctx))
	})
}

func (hmw OtelHmw) match(req *http.Request) (*middleware.MatchedRoute, bool) {
	route, matched := hmw.mctx.LookupRoute(req)

	return route, matched && hmw.routeTester(route)
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

// HttpTripperwareStack chains HTTP client middleware. The first Push runs first.
type HttpTripperwareStack = tricks.MiddlewareStack[http.RoundTripper]

// BlindLoggerHttpTripperwareOpt configures HttpBlindLoggerTripperware. No options are defined yet.
type BlindLoggerHttpTripperwareOpt = tricks.MutableOption[any]

// HttpBlindLoggerTripperware records the dumped request and response on the current span.
func HttpBlindLoggerTripperware(options ...BlindLoggerHttpTripperwareOpt) tricks.Middleware[http.RoundTripper] {
	return func(next http.RoundTripper) http.RoundTripper {
		return HttpRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			span := trace.SpanFromContext(req.Context())
			attrs := make([]attribute.KeyValue, 0)
			defer func() {
				span.AddEvent("http request", trace.WithAttributes(attrs...))
			}()

			bb, _ := httputil.DumpRequestOut(req, true)
			attrs = append(attrs, oaHttpRequest.String(string(bb)))

			rsp, err := next.RoundTrip(req)
			span.RecordError(err)
			if rsp != nil {
				bb, _ = httputil.DumpResponse(rsp, true)
				attrs = append(attrs, oaHttpResponse.String(string(bb)))
			}

			return rsp, err
		})
	}
}

// HttpErrorMapperTripperware passes the next round trip's response and error through mapper.
// It panics when mapper is nil.
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

// HttpRoundTripperFunc adapts a function to http.RoundTripper.
type HttpRoundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip calls fn.
func (fn HttpRoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
