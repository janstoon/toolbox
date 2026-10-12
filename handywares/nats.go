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

// NatsMsgHandler handles one NATS message.
type NatsMsgHandler func(ctx context.Context, msg *nats.Msg) error

// NatsMiddlewareStack chains NATS handlers. The first Push runs first.
type NatsMiddlewareStack = tricks.MiddlewareStack[NatsMsgHandler]

// PanicRecoverNatsMiddlewareOpt configures NatsPanicRecoverMiddleware. No options are defined yet.
type PanicRecoverNatsMiddlewareOpt = tricks.MutableOption[any]

// NatsPanicRecoverMiddleware recovers a panic in the next handler and records it on the current span.
// The recovered call returns a nil error.
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

// CompensatorNatsMiddlewareOpt configures NatsCompensatorMiddleware. No options are defined yet.
type CompensatorNatsMiddlewareOpt = tricks.MutableOption[any]

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

// OtelNmw is the OpenTelemetry NATS middleware built by NatsOpenTelemetryMiddleware.
type OtelNmw struct {
	tracer trace.Tracer

	namePrefix string
}

// OpenTelemetryNatsMiddlewareOpt configures NatsOpenTelemetryMiddleware.
type OpenTelemetryNatsMiddlewareOpt = tricks.Option[OtelNmw]

// OtelNatsSpanNamePrefix sets the text written before the subject in the span name.
func OtelNatsSpanNamePrefix(prefix string) OpenTelemetryNatsMiddlewareOpt {
	return tricks.ImmutableOption[OtelNmw](func(nmw OtelNmw) OtelNmw {
		nmw.namePrefix = prefix

		return nmw
	})
}

// NatsOpenTelemetryMiddleware starts a consumer span named after the message subject.
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
