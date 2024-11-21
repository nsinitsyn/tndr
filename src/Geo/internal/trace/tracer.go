package trace

import (
	"context"
	"runtime"
	"strings"
	"tinder-geo/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var tracingEnabled bool = false

func InitTracer(config config.TracingConfig, serviceConfig config.ServiceConfig) error {
	tracingEnabled = config.Enabled

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

func GetStackTraceData() (file string, line int, function string) {
	return getStackTraceDataWithSkip(2)
}

func StartNewSpanWithCurrentFunctionName(ctx context.Context, tracer trace.Tracer, kv ...attribute.KeyValue) (context.Context, trace.Span) {
	if !tracingEnabled {
		return ctx, noop.Span{}
	}

	_, _, function := getStackTraceDataWithSkip(3)
	function = function[strings.LastIndex(function, "/")+1:]

	return tracer.Start(
		ctx,
		function,
		trace.WithAttributes(kv...),
	)
}

func AddAttributesToCurrentSpan(ctx context.Context, kv ...attribute.KeyValue) {
	if !tracingEnabled {
		return
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(kv...)
}

func SetErrorForCurrentSpan(ctx context.Context, err error) {
	if !tracingEnabled {
		return
	}

	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
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

func getStackTraceDataWithSkip(skip int) (file string, line int, function string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.File, frame.Line, frame.Function
}
