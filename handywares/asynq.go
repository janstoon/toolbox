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

type AsynqMiddlewareStack asynq.MiddlewareFunc

type PanicRecoverAsynqMiddlewareOpt = tricks.InPlaceOption[any]

func (stk *AsynqMiddlewareStack) PushPanicRecover(options ...PanicRecoverAsynqMiddlewareOpt) *AsynqMiddlewareStack {
	return stk.Push(func(next asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("paniced %+v\n", r)
					debug.PrintStack()
				}
			}()

			return next.ProcessTask(ctx, task)
		})
	})
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

func (stk *AsynqMiddlewareStack) PushOpenTelemetry(
	tracer trace.Tracer, options ...OpenTelemetryAsynqMiddlewareOpt,
) *AsynqMiddlewareStack {
	amw := &OtelAmw{
		tracer: tracer,
	}
	amw = tricks.ApplyOptions(amw, options...)

	return stk.Push(amw.builder)
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

func asynqPassthroughBuilder(handler asynq.Handler) asynq.Handler { return handler }

func (stk *AsynqMiddlewareStack) Push(mw asynq.MiddlewareFunc) *AsynqMiddlewareStack {
	current := *stk
	if current == nil {
		current = asynqPassthroughBuilder
	}

	*stk = func(next asynq.Handler) asynq.Handler {
		return current(mw(next))
	}

	return stk
}

func (stk *AsynqMiddlewareStack) Cast() asynq.MiddlewareFunc {
	return asynq.MiddlewareFunc(*stk)
}
