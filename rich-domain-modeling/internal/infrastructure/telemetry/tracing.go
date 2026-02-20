package telemetry

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// SetupTracing configures a global OpenTelemetry tracer provider.
func SetupTracing(ctx context.Context, serviceName string) (func(context.Context) error, error) {
	exporterType := strings.ToLower(strings.TrimSpace(os.Getenv("TRACE_EXPORTER")))
	if exporterType == "" {
		exporterType = "stdout"
	}

	var (
		exporter tracesdk.SpanExporter
		err      error
	)

	switch exporterType {
	case "stdout":
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, err
		}
	case "zipkin":
		endpoint := strings.TrimSpace(os.Getenv("ZIPKIN_ENDPOINT"))
		if endpoint == "" {
			endpoint = "http://localhost:9411/api/v2/spans"
		}
		exporter = NewZipkinExporter(endpoint, serviceName)
	default:
		return nil, fmt.Errorf("unsupported TRACE_EXPORTER %q", exporterType)
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(attribute.String("service.name", serviceName)),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(res),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}
