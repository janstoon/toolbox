package handywares_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/handywares"
)

func TestHttpHelpers(t *testing.T) {
	always := handywares.PassthroughHttpRouteTester(true)
	never := handywares.PassthroughHttpRouteTester(false)
	assert.True(t, always(nil))
	assert.False(t, never(nil))
	assert.True(t, handywares.CombineHttpRouteTesters()(nil))
	assert.True(t, handywares.CombineHttpRouteTesters(always)(nil))
	assert.False(t, handywares.CombineHttpRouteTesters(always, never)(nil))

	handler := handywares.HttpCrossOriginResourceSharingPolicyMiddleware(
		handywares.CorsAllowOrigins("https://example.com"),
		handywares.CorsAllowMethods(http.MethodGet),
		handywares.CorsAllowHeaders("X-Request-Id"),
		handywares.CorsDebug(false),
	)(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodOptions, "http://app.example/", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, "https://example.com", rec.Header().Get("Access-Control-Allow-Origin"))

	var logged handywares.HttpMiddlewareStack
	logged = logged.Push(handywares.HttpBlindLoggerMiddleware(nil))
	rec = httptest.NewRecorder()
	logged(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "http://app.example/health", nil))
	assert.Equal(t, http.StatusNoContent, rec.Code)

	next := handywares.HttpRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(strings.NewReader("ok")),
			Header:     make(http.Header),
			Request:    req,
		}, nil
	})
	var trip handywares.HttpTripperwareStack
	trip = trip.Push(handywares.HttpBlindLoggerTripperware()).
		Push(handywares.HttpErrorMapperTripperware(handywares.HttpToBricksErrorMapper))
	rsp, err := trip(next).RoundTrip(httptest.NewRequest(http.MethodPost, "http://app.example/items", strings.NewReader("body")))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rsp.StatusCode)
	require.NoError(t, rsp.Body.Close())

	assert.Panics(t, func() {
		handywares.HttpErrorMapperTripperware(nil)
	})

	tracer := noop.NewTracerProvider().Tracer("test")
	assert.NotNil(t, handywares.HttpOpenTelemetryMiddleware(
		tracer,
		nil,
		handywares.OtelHttpSpanNamePrefix("api"),
		handywares.OtelHttpRouteTester(handywares.PassthroughHttpRouteTester(true)),
		handywares.OtelHttpOperationIdException("health"),
	))
}

func TestGrpcMiddleware(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("test")
	info := &grpc.UnaryServerInfo{FullMethod: "/pkg.Svc/Get"}
	var stack handywares.GrpcUnaryServerMiddlewareStack
	stack = stack.
		Push(handywares.GrpcPanicRecoverMiddleware()).
		Push(handywares.GrpcOpenTelemetryMiddleware(tracer, handywares.OtelGrpcSpanNamePrefix("grpc"))).
		Push(handywares.GrpcUnaryServerErrorMapperMiddleware(handywares.BricksToGrpcErrorMapper))

	_, err := stack(handywares.GrpcUnaryServerInvokeHandlerInterceptor)(
		context.Background(), nil, info,
		func(context.Context, any) (any, error) {
			panic("boom")
		},
	)
	assert.NoError(t, err)

	_, err = stack(handywares.GrpcUnaryServerInvokeHandlerInterceptor)(
		context.Background(), nil, info,
		func(context.Context, any) (any, error) {
			return nil, bricks.ErrNotFound
		},
	)
	assert.Error(t, err)

	assert.Panics(t, func() {
		handywares.GrpcUnaryServerErrorMapperMiddleware(nil)
	})

	var clients handywares.GrpcUnaryClientMiddlewareStack
	clients = clients.
		Push(handywares.GrpcUnaryClientErrorMapperMiddleware(handywares.GrpcToBricksErrorMapper))
	err = clients(handywares.GrpcUnaryClientInvokerInterceptor)(
		context.Background(), "/pkg.Svc/Get", nil, nil, nil,
		func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
			return nil
		},
	)
	assert.NoError(t, err)
	assert.Panics(t, func() {
		handywares.GrpcUnaryClientErrorMapperMiddleware(nil)
	})
}

func TestAsynqMiddleware(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("test")
	task := asynq.NewTask("email:send", []byte("x"))

	var stack handywares.AsynqMiddlewareStack
	stack = stack.
		Push(handywares.AsynqPanicRecoverMiddleware()).
		Push(handywares.AsynqOpenTelemetryMiddleware(tracer, handywares.OtelAsynqSpanNamePrefix("jobs"))).
		Push(handywares.AsynqErrorDecoratorMiddleware()).
		Push(handywares.AsynqCompensatorMiddleware())

	err := stack(asynq.HandlerFunc(func(context.Context, *asynq.Task) error {
		panic("boom")
	})).ProcessTask(context.Background(), task)
	assert.NoError(t, err)

	err = stack(asynq.HandlerFunc(func(context.Context, *asynq.Task) error {
		return errors.New("nope")
	})).ProcessTask(context.Background(), task)
	require.ErrorIs(t, err, asynq.SkipRetry)

	compensated := false
	err = stack(asynq.HandlerFunc(func(context.Context, *asynq.Task) error {
		return bricks.CompensatorAsError(bricks.CompensatorFunc(func(context.Context, error) error {
			compensated = true

			return errors.New("undone")
		}), errors.New("failed"))
	})).ProcessTask(context.Background(), task)
	require.Error(t, err)
	assert.True(t, compensated)

	err = handywares.AsynqErrorDecoratorMiddleware()(asynq.HandlerFunc(func(context.Context, *asynq.Task) error {
		return bricks.ErrUnavailable
	})).ProcessTask(context.Background(), task)
	require.ErrorIs(t, err, bricks.ErrUnavailable)
	assert.NotErrorIs(t, err, asynq.SkipRetry)
}

func TestMessageAndNatsMiddleware(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("test")

	var stack handywares.MsgMiddlewareStack[string]
	stack = stack.Push(handywares.MsgPanicRecoverMiddleware[string]()).
		Push(handywares.MsgOpenTelemetryMiddleware(
			tracer,
			handywares.OtelMsgSpanNamePrefix[string]("msg"),
			handywares.OtelMsgSubjectExporter(func(msg string) string { return msg }),
		)).
		Push(handywares.MsgCompensatorMiddleware[string]())

	err := stack(func(context.Context, string) error {
		panic("boom")
	})(context.Background(), "jobs.email")
	assert.NoError(t, err)

	compensated := false
	err = stack(func(context.Context, string) error {
		return bricks.CompensatorAsError(bricks.CompensatorFunc(func(context.Context, error) error {
			compensated = true

			return nil
		}), errors.New("failed"))
	})(context.Background(), "jobs.email")
	require.Error(t, err)
	assert.True(t, compensated)

	err = stack(func(context.Context, string) error {
		return bricks.ErrUnavailable
	})(context.Background(), "jobs.email")
	require.ErrorIs(t, err, bricks.ErrUnavailable)

	var natsStack handywares.NatsMiddlewareStack
	natsStack = natsStack.Push(handywares.NatsPanicRecoverMiddleware()).
		Push(handywares.NatsOpenTelemetryMiddleware(tracer, handywares.OtelNatsSpanNamePrefix("nats"))).
		Push(handywares.NatsCompensatorMiddleware())

	msg := &nats.Msg{Subject: "jobs.email"}
	err = natsStack(func(context.Context, *nats.Msg) error {
		panic("boom")
	})(context.Background(), msg)
	assert.NoError(t, err)

	err = natsStack(func(context.Context, *nats.Msg) error {
		return nil
	})(context.Background(), msg)
	assert.NoError(t, err)
}

func ExamplePassthroughHttpRouteTester() {
	tester := handywares.CombineHttpRouteTesters(
		handywares.PassthroughHttpRouteTester(true),
		handywares.PassthroughHttpRouteTester(true),
	)
	fmt.Println(tester(nil))

	// Output:
	// true
}
