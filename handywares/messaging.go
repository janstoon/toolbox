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
)

// MsgHandler handles one message of type M.
type MsgHandler[M any] func(ctx context.Context, msg M) error

// MsgMiddlewareStack chains message handlers. The first Push runs first.
type MsgMiddlewareStack[M any] tricks.MiddlewareStack[MsgHandler[M]]

// Push is a hack to provide underlying tricks.MiddlewareStack.Push
// todo: make MsgMiddlewareStack[M any] type alias of tricks.MiddlewareStack[MsgHandler[M]] go 1.24 and remove Push
func (stk MsgMiddlewareStack[M]) Push(mw tricks.Middleware[MsgHandler[M]]) MsgMiddlewareStack[M] {
	return MsgMiddlewareStack[M](tricks.MiddlewareStack[MsgHandler[M]](stk).Push(mw))
}

// PanicRecoverMsgMiddlewareOpt configures MsgPanicRecoverMiddleware. No options are defined yet.
type PanicRecoverMsgMiddlewareOpt = tricks.MutableOption[any]

// MsgPanicRecoverMiddleware recovers a panic in the next handler and records it on the current span.
// The recovered call returns a nil error.
func MsgPanicRecoverMiddleware[M any](options ...PanicRecoverMsgMiddlewareOpt) tricks.Middleware[MsgHandler[M]] {
	return func(next MsgHandler[M]) MsgHandler[M] {
		return func(ctx context.Context, msg M) error {
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

// CompensatorMsgMiddlewareOpt configures MsgCompensatorMiddleware. No options are defined yet.
type CompensatorMsgMiddlewareOpt = tricks.MutableOption[any]

// MsgCompensatorMiddleware tries to compensate the task if error is not bricks.ErrRetryable.
// It searches for a bricks.Compensator in returned error by msg handler and runs the first one.
func MsgCompensatorMiddleware[M any](options ...CompensatorMsgMiddlewareOpt) tricks.Middleware[MsgHandler[M]] {
	return func(next MsgHandler[M]) MsgHandler[M] {
		return func(ctx context.Context, msg M) error {
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

// OtelMmw is the OpenTelemetry message middleware built by MsgOpenTelemetryMiddleware.
type OtelMmw[M any] struct {
	tracer trace.Tracer

	namePrefix   string
	subjExporter func(msg M) string
}

// OpenTelemetryMsgMiddlewareOpt configures MsgOpenTelemetryMiddleware.
type OpenTelemetryMsgMiddlewareOpt[M any] tricks.Option[OtelMmw[M]]

// OtelMsgSpanNamePrefix sets the text written before the subject in the span name.
func OtelMsgSpanNamePrefix[M any](prefix string) OpenTelemetryMsgMiddlewareOpt[M] {
	return tricks.ImmutableOption[OtelMmw[M]](func(nmw OtelMmw[M]) OtelMmw[M] {
		nmw.namePrefix = prefix

		return nmw
	})
}

// OtelMsgSubjectExporter sets the function that names the span from the message.
// The default name is "unknown".
func OtelMsgSubjectExporter[M any](exp func(msg M) string) OpenTelemetryMsgMiddlewareOpt[M] {
	return tricks.ImmutableOption[OtelMmw[M]](func(mmw OtelMmw[M]) OtelMmw[M] {
		mmw.subjExporter = exp

		return mmw
	})
}

// MsgOpenTelemetryMiddleware starts a consumer span named by the subject exporter.
func MsgOpenTelemetryMiddleware[M any](
	tracer trace.Tracer, options ...OpenTelemetryMsgMiddlewareOpt[M],
) tricks.Middleware[MsgHandler[M]] {
	amw := &OtelMmw[M]{
		tracer: tracer,

		subjExporter: func(msg M) string {
			return "unknown"
		},
	}
	amw = tricks.ApplyOptions(amw, tricks.SliceMap[[]OpenTelemetryMsgMiddlewareOpt[M], []tricks.Option[OtelMmw[M]]](
		options,
		func(src OpenTelemetryMsgMiddlewareOpt[M]) tricks.Option[OtelMmw[M]] {
			return src
		},
	)...)

	return amw.builder
}

func (mmw OtelMmw[M]) builder(next MsgHandler[M]) MsgHandler[M] {
	return func(ctx context.Context, msg M) error {
		var span trace.Span
		ctx, span = mmw.tracer.Start(ctx, mmw.spanName(mmw.subjExporter(msg)), trace.WithSpanKind(trace.SpanKindConsumer))
		defer span.End()

		return next(ctx, msg)
	}
}

func (mmw OtelMmw[M]) spanName(subject string) string {
	sb := strings.Builder{}
	if len(strings.TrimSpace(mmw.namePrefix)) > 0 {
		sb.WriteString(mmw.namePrefix)
		sb.WriteRune('/')
	}

	sb.WriteString(subject)

	return sb.String()
}
