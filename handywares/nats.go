package handywares

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/tricks"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/trace"
)

type NatsMsgHandler func(ctx context.Context, msg *nats.Msg) error

type NatsMiddlewareStack = tricks.MiddlewareStack[NatsMsgHandler]

type PanicRecoverNatsMiddlewareOpt = tricks.InPlaceOption[any]

func NatsPanicRecoverMiddleware(options ...PanicRecoverAsynqMiddlewareOpt) tricks.Middleware[NatsMsgHandler] {
	return func(next NatsMsgHandler) NatsMsgHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			defer func() {
				if r := recover(); r != nil {
					span := trace.SpanFromContext(ctx)
					span.AddEvent("panic recovered", trace.WithAttributes(
						oaPanicValue.String(fmt.Sprintf("%+v", r)),
						oaDebugStack.String(string(debug.Stack())),
					))
				}
			}()

			return next(ctx, msg)
		}
	}
}

type CompensatorNatsMiddlewareOpt = tricks.InPlaceOption[any]

// NatsCompensatorMiddleware tries to compensate the task if error is not bricks.ErrRetryable.
// It searches for a bricks.Compensator in returned error by msg handler and runs the first one.
func NatsCompensatorMiddleware(options ...CompensatorNatsMiddlewareOpt) tricks.Middleware[NatsMsgHandler] {
	return func(next NatsMsgHandler) NatsMsgHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			err := next(ctx, msg)
			if err == nil {
				return nil
			}

			if !errors.Is(err, bricks.ErrRetryable) {
				var c bricks.Compensator
				if errors.As(err, &c) {
					err = errors.Join(err, c.Compensate(ctx, err))
				}
			}

			return err
		}
	}
}

type OtelNmw struct {
	tracer trace.Tracer

	namePrefix string
}

type OpenTelemetryNatsMiddlewareOpt = tricks.Option[OtelNmw]

func OtelNatsSpanNamePrefix(prefix string) OpenTelemetryNatsMiddlewareOpt {
	return tricks.OutOfPlaceOption[OtelNmw](func(nmw OtelNmw) OtelNmw {
		nmw.namePrefix = prefix

		return nmw
	})
}

func NatsOpenTelemetryMiddleware(
	tracer trace.Tracer, options ...OpenTelemetryNatsMiddlewareOpt,
) tricks.Middleware[NatsMsgHandler] {
	amw := &OtelNmw{
		tracer: tracer,
	}
	amw = tricks.ApplyOptions(amw, options...)

	return amw.builder
}

func (nmw OtelNmw) builder(next NatsMsgHandler) NatsMsgHandler {
	return func(ctx context.Context, msg *nats.Msg) error {
		var span trace.Span
		ctx, span = nmw.tracer.Start(ctx, nmw.spanName(msg.Subject), trace.WithSpanKind(trace.SpanKindConsumer))
		defer span.End()

		return next(ctx, msg)
	}
}

func (nmw OtelNmw) spanName(subject string) string {
	sb := strings.Builder{}
	if len(strings.TrimSpace(nmw.namePrefix)) > 0 {
		sb.WriteString(nmw.namePrefix)
		sb.WriteRune('/')
	}

	sb.WriteString(subject)

	return sb.String()
}
