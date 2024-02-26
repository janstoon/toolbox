package std

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.24.0"

	"github.com/janstoon/toolbox/kareless"
)

const (
	TelemetryCollectorOtlpHostSettingKey  = "telemetry.collector.otlp.host"
	TelemetryCollectorOtlpPortsSettingKey = "telemetry.collector.otlp.ports"
)

func otlpGrpcEndpoint(ss *kareless.Settings) string {
	return otlpEndpoint("grpc", ss)
}

func otlpHttpEndpoint(ss *kareless.Settings) string {
	return otlpEndpoint("http", ss)
}

func otlpEndpoint(protocol string, ss *kareless.Settings) string {
	var (
		host = strings.TrimSpace(ss.GetString(TelemetryCollectorOtlpHostSettingKey))
		port = ss.GetInt(fmt.Sprintf("%s.%s", TelemetryCollectorOtlpPortsSettingKey, protocol))
	)

	return fmt.Sprintf("%s:%d", host, port)
}

type otlp struct {
	serviceName string
	attrs       []attribute.KeyValue

	client otlptrace.Client
}

func OpenTelemetryDriverConstructor(serviceName string, attrs ...attribute.KeyValue) kareless.DriverConstructor {
	return func(ss *kareless.Settings, _ *kareless.InstrumentBank, _ []kareless.Application) kareless.Driver {
		client := otlptracegrpc.NewClient(
			otlptracegrpc.WithEndpoint(otlpGrpcEndpoint(ss)),
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithRetry(
				otlptracegrpc.RetryConfig{
					Enabled:         true,
					InitialInterval: time.Second,
					MaxInterval:     10 * time.Second,
					MaxElapsedTime:  time.Minute,
				},
			),
		)

		return otlp{
			serviceName: serviceName,
			attrs:       attrs,

			client: client,
		}
	}
}

func (d otlp) Run(ctx context.Context) error {
	traceExp, err := otlptrace.New(ctx, d.client)
	if err != nil {
		return err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(traceExp),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				append(d.attrs, semconv.ServiceNameKey.String(d.serviceName))...,
			),
		),
	)
	otel.SetTracerProvider(tp)

	promRegistry := prom.NewRegistry()
	promRegistry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	metricExp, err := prometheus.New(
		prometheus.WithRegisterer(promRegistry),
	)
	if err != nil {
		return err
	}

	mp := metric.NewMeterProvider(metric.WithReader(metricExp))
	otel.SetMeterProvider(mp)

	<-ctx.Done()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		_ = tp.Shutdown(context.Background())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		_ = mp.Shutdown(context.Background())
	}()

	wg.Wait()

	return nil
}
