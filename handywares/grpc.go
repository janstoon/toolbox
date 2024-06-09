package handywares

import (
	"context"
	"errors"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/tricks"
	"google.golang.org/grpc"
)

type GrpcUnaryServerMiddlewareStack = tricks.MiddlewareStack[grpc.UnaryServerInterceptor]

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
