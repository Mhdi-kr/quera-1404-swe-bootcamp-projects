package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type ZipkinExporter struct {
	client      *http.Client
	endpoint    string
	serviceName string
}

func NewZipkinExporter(endpoint, serviceName string) *ZipkinExporter {
	return &ZipkinExporter{
		client:      &http.Client{Timeout: 5 * time.Second},
		endpoint:    endpoint,
		serviceName: serviceName,
	}
}

func (e *ZipkinExporter) ExportSpans(ctx context.Context, spans []tracesdk.ReadOnlySpan) error {
	if len(spans) == 0 {
		return nil
	}

	zipkinSpans := make([]zipkinSpan, 0, len(spans))
	for _, span := range spans {
		zipkinSpans = append(zipkinSpans, e.mapSpan(span))
	}

	payload, err := json.Marshal(zipkinSpans)
	if err != nil {
		return fmt.Errorf("marshal zipkin payload: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, e.endpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("build zipkin request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := e.client.Do(request)
	if err != nil {
		return fmt.Errorf("send zipkin spans: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 2048))
		return fmt.Errorf("zipkin endpoint returned %d: %s", response.StatusCode, strings.TrimSpace(string(body)))
	}

	return nil
}

func (e *ZipkinExporter) Shutdown(_ context.Context) error {
	return nil
}

func (e *ZipkinExporter) mapSpan(span tracesdk.ReadOnlySpan) zipkinSpan {
	spanContext := span.SpanContext()
	parent := span.Parent()

	mapped := zipkinSpan{
		TraceID:   spanContext.TraceID().String(),
		ID:        spanContext.SpanID().String(),
		Name:      span.Name(),
		Timestamp: span.StartTime().UnixMicro(),
		Duration:  span.EndTime().Sub(span.StartTime()).Microseconds(),
		LocalEndpoint: zipkinEndpoint{
			ServiceName: e.serviceName,
		},
		Tags: make(map[string]string),
	}

	if parent.IsValid() {
		mapped.ParentID = parent.SpanID().String()
	}

	for _, attribute := range span.Attributes() {
		mapped.Tags[string(attribute.Key)] = attributeValueToString(attribute.Value)
	}

	status := span.Status()
	if status.Code != codes.Unset {
		mapped.Tags["otel.status_code"] = status.Code.String()
	}
	if status.Description != "" {
		mapped.Tags["otel.status_description"] = status.Description
	}

	return mapped
}

func attributeValueToString(value attribute.Value) string {
	switch value.Type() {
	case attribute.BOOL:
		return strconv.FormatBool(value.AsBool())
	case attribute.INT64:
		return strconv.FormatInt(value.AsInt64(), 10)
	case attribute.FLOAT64:
		return strconv.FormatFloat(value.AsFloat64(), 'f', -1, 64)
	case attribute.STRING:
		return value.AsString()
	case attribute.BOOLSLICE:
		return joinSlice(value.AsBoolSlice(), strconv.FormatBool)
	case attribute.INT64SLICE:
		return joinSlice(value.AsInt64Slice(), func(v int64) string { return strconv.FormatInt(v, 10) })
	case attribute.FLOAT64SLICE:
		return joinSlice(value.AsFloat64Slice(), func(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) })
	case attribute.STRINGSLICE:
		return strings.Join(value.AsStringSlice(), ",")
	default:
		return ""
	}
}

func joinSlice[T any](values []T, toString func(T) string) string {
	if len(values) == 0 {
		return ""
	}

	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, toString(value))
	}

	return strings.Join(result, ",")
}

type zipkinEndpoint struct {
	ServiceName string `json:"serviceName,omitempty"`
}

type zipkinSpan struct {
	TraceID       string            `json:"traceId"`
	ID            string            `json:"id"`
	ParentID      string            `json:"parentId,omitempty"`
	Name          string            `json:"name"`
	Timestamp     int64             `json:"timestamp"`
	Duration      int64             `json:"duration"`
	LocalEndpoint zipkinEndpoint    `json:"localEndpoint,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
}
