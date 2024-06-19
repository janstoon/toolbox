package handywares

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/tricks"
	"go.opentelemetry.io/otel/trace"
)

type AsynqMiddlewareStack = tricks.MiddlewareStack[asynq.Handler]

type ErrorDecoratorAsynqMiddlewareOpt = tricks.InPlaceOption[any]

// AsynqErrorDecoratorMiddleware skip asynq retry except that error is bricks.ErrRetryable
func AsynqErrorDecoratorMiddleware(options ...ErrorDecoratorAsynqMiddlewareOpt) tricks.Middleware[asynq.Handler] {
	return func(next asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
			err := next.ProcessTask(ctx, task)
			if err == nil {
				return nil
			}

			if errors.Is(err, bricks.ErrRetryable) {
				return err
			}

			return errors.Join(asynq.SkipRetry, err)
		})
	}
}

type PanicRecoverAsynqMiddlewareOpt = tricks.InPlaceOption[any]

func AsynqPanicRecoverMiddleware(options ...PanicRecoverAsynqMiddlewareOpt) tricks.Middleware[asynq.Handler] {
	return func(next asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
			defer func() {
				if r := recover(); r != nil {
					span := trace.SpanFromContext(ctx)
					span.AddEvent("panic recovered", trace.WithAttributes(
						oaPanicValue.String(fmt.Sprintf("%+v", r)),
						oaDebugStack.String(string(debug.Stack())),
					))
				}
			}()

			return next.ProcessTask(ctx, task)
		})
	}
}

type CompensatorAsynqMiddlewareOpt = tricks.InPlaceOption[any]

// AsynqCompensatorMiddleware tries to compensate the task if max retries reached or error is not bricks.ErrRetryable.
// It searches for a bricks.Compensator in returned error by task handler and runs the first one.
func AsynqCompensatorMiddleware(options ...CompensatorAsynqMiddlewareOpt) tricks.Middleware[asynq.Handler] {
	return func(next asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
			err := next.ProcessTask(ctx, task)
			if err == nil {
				return nil
			}

			retries, _ := asynq.GetRetryCount(ctx)
			maxRetry, _ := asynq.GetMaxRetry(ctx)
			if retries >= maxRetry || !errors.Is(err, bricks.ErrRetryable) {
				var c bricks.Compensator
				if errors.As(err, &c) {
					err = errors.Join(err, c.Compensate(ctx, err))
				}
			}

			return err
		})
	}
}

type OtelAmw struct {
	tracer trace.Tracer

	namePrefix string
}

type OpenTelemetryAsynqMiddlewareOpt = tricks.Option[OtelAmw]

func OtelAsynqSpanNamePrefix(prefix string) OpenTelemetryAsynqMiddlewareOpt {
	return tricks.OutOfPlaceOption[OtelAmw](func(amw OtelAmw) OtelAmw {
		amw.namePrefix = prefix

		return amw
	})
}

func AsynqOpenTelemetryMiddleware(
	tracer trace.Tracer, options ...OpenTelemetryAsynqMiddlewareOpt,
) tricks.Middleware[asynq.Handler] {
	amw := &OtelAmw{
		tracer: tracer,
	}
	amw = tricks.ApplyOptions(amw, options...)

	return amw.builder
}

func (amw OtelAmw) builder(next asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
		var span trace.Span
		ctx, span = amw.tracer.Start(ctx, amw.spanName(task.Type()), trace.WithSpanKind(trace.SpanKindConsumer))
		defer span.End()

		return next.ProcessTask(ctx, task)
	})
}

func (amw OtelAmw) spanName(opId string) string {
	sb := strings.Builder{}
	if len(strings.TrimSpace(amw.namePrefix)) > 0 {
		sb.WriteString(amw.namePrefix)
		sb.WriteRune('/')
	}

	sb.WriteString(opId)

	return sb.String()
}
