package trace

import (
	"context"
	"tinder-geo/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func InitTracer(config config.TracingConfig, serviceConfig config.ServiceConfig) error {
	traceClient := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(config.Endpoint),
		otlptracehttp.WithInsecure(),
	)
	exporter, err := otlptrace.New(context.Background(), traceClient)
	if err != nil {
		return err
	}

	provider, err := newTraceProvider(exporter, serviceConfig)
	if err != nil {
		return err
	}

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return nil
}

func newTraceProvider(exporter *otlptrace.Exporter, serviceConfig config.ServiceConfig) (*sdk_trace.TracerProvider, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.DeploymentEnvironmentKey.String(serviceConfig.Env),
			semconv.ServiceNameKey.String(serviceConfig.Name),
			semconv.ServiceVersionKey.String(serviceConfig.Version),
			semconv.ServiceInstanceIDKey.String(serviceConfig.InstanceID),
		),
	)
	if err != nil {
		return nil, err
	}

	provider := sdk_trace.NewTracerProvider(
		sdk_trace.WithBatcher(exporter),
		sdk_trace.WithResource(r),
	)

	return provider, nil
}
