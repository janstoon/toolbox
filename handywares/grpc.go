package handywares

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/tricks"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type GrpcUnaryServerMiddlewareStack = tricks.MiddlewareStack[grpc.UnaryServerInterceptor]

type PanicRecoverGrpcMiddlewareOpt = tricks.InPlaceOption[any]

func GrpcPanicRecoverMiddleware(
	options ...PanicRecoverGrpcMiddlewareOpt,
) tricks.Middleware[grpc.UnaryServerInterceptor] {
	return func(next grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
		return func(
			ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
		) (resp any, err error) {
			defer func() {
				if r := recover(); r != nil {
					span := trace.SpanFromContext(ctx)
					span.AddEvent("panic recovered", trace.WithAttributes(
						oaPanicValue.String(fmt.Sprintf("%+v", r)),
						oaDebugStack.String(string(debug.Stack())),
					))
				}
			}()

			return next(ctx, req, info, handler)
		}
	}
}

type OtelGmw struct {
	tracer trace.Tracer

	namePrefix string
}

type OpenTelemetryGrpcMiddlewareOpt = tricks.Option[OtelGmw]

func OtelGrpcSpanNamePrefix(prefix string) OpenTelemetryGrpcMiddlewareOpt {
	return tricks.OutOfPlaceOption[OtelGmw](func(gmw OtelGmw) OtelGmw {
		gmw.namePrefix = prefix

		return gmw
	})
}

func GrpcOpenTelemetryMiddleware(
	tracer trace.Tracer, options ...OpenTelemetryGrpcMiddlewareOpt,
) tricks.Middleware[grpc.UnaryServerInterceptor] {
	gmw := &OtelGmw{
		tracer: tracer,
	}
	gmw = tricks.ApplyOptions(gmw, options...)

	return gmw.builder
}

func (gmw OtelGmw) builder(next grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		var span trace.Span
		ctx, span = gmw.tracer.Start(ctx, gmw.spanName(info.FullMethod), trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		return next(ctx, req, info, handler)
	}
}

func (gmw OtelGmw) spanName(opId string) string {
	sb := strings.Builder{}
	if len(strings.TrimSpace(gmw.namePrefix)) > 0 {
		sb.WriteString(gmw.namePrefix)
		sb.WriteRune('/')
	}

	sb.WriteString(opId)

	return sb.String()
}

func GrpcUnaryServerInvokeHandlerInterceptor(
	ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (any, error) {
	return handler(ctx, req)
}

func GrpcUnaryServerErrorMapperMiddleware(mapper func(error) error) tricks.Middleware[grpc.UnaryServerInterceptor] {
	if mapper == nil {
		panic(errors.Join(bricks.ErrInvalidArgument, errors.New("empty error mapper")))
	}

	return func(next grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
		return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
			rsp, err := next(ctx, req, info, handler)

			return rsp, mapper(err)
		}
	}
}

type GrpcUnaryClientMiddlewareStack = tricks.MiddlewareStack[grpc.UnaryClientInterceptor]

func GrpcUnaryClientInvokerInterceptor(
	ctx context.Context, method string, req, reply any,
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {
	return invoker(ctx, method, req, reply, cc, opts...)
}

func GrpcUnaryClientErrorMapperMiddleware(mapper func(error) error) tricks.Middleware[grpc.UnaryClientInterceptor] {
	if mapper == nil {
		panic(errors.Join(bricks.ErrInvalidArgument, errors.New("empty error mapper")))
	}

	return func(next grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
		return func(
			ctx context.Context, method string, req, reply any,
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
		) error {
			err := next(ctx, method, req, reply, cc, invoker, opts...)

			return mapper(err)
		}
	}
}
