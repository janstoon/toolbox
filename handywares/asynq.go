package handywares

import (
	"context"
	"log"
	"runtime/debug"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/janstoon/toolbox/tricks"
	"go.opentelemetry.io/otel/trace"
)

type AsynqMiddlewareStack = tricks.MiddlewareStack[asynq.Handler]

type PanicRecoverAsynqMiddlewareOpt = tricks.InPlaceOption[any]

func AsynqPanicRecoverMiddleware(options ...PanicRecoverAsynqMiddlewareOpt) tricks.Middleware[asynq.Handler] {
	return func(next asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("paniced %+v\n", r)
					debug.PrintStack()
				}
			}()

			return next.ProcessTask(ctx, task)
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
		ctx, span = amw.tracer.Start(ctx, amw.spanName(task.Type()))
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
