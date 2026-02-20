package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"

	"richdomainmodeling/internal/infrastructure/eventbus"
	"richdomainmodeling/internal/infrastructure/persistence/inmemory"
	"richdomainmodeling/internal/interfaces/httpapi"
	kitchenapp "richdomainmodeling/internal/kitchen/application"
	reservationapp "richdomainmodeling/internal/reservation/application"
	theaterapp "richdomainmodeling/internal/theater/application"
	theaterdomain "richdomainmodeling/internal/theater/domain"
)

func TestVIPSeatPurchaseIssuesVoucherViaHTTP(t *testing.T) {
	t.Parallel()

	server := newTestServer(t)
	handler := server.Handler()

	startsAt := time.Date(2026, 2, 20, 20, 0, 0, 0, time.UTC)

	schedulePayload := map[string]any{
		"show_id":   "show-http-1",
		"title":     "HTTP VIP Night",
		"starts_at": startsAt.Format(time.RFC3339),
		"seats": []map[string]any{
			{"number": "V1", "tier": "vip", "price_cents": 9000},
		},
	}
	doJSONRequest(t, handler, http.MethodPost, "/shows", schedulePayload, http.StatusCreated, nil)

	purchasePayload := map[string]any{
		"customer_id": "member-77",
		"seat_number": "V1",
	}

	var ticket struct {
		IncludesFreeCoffee bool `json:"includes_free_coffee"`
	}
	doJSONRequest(t, handler, http.MethodPost, "/shows/show-http-1/purchase", purchasePayload, http.StatusCreated, &ticket)
	if !ticket.IncludesFreeCoffee {
		t.Fatalf("expected vip ticket to include free coffee")
	}

	var vouchers []struct {
		ID string `json:"id"`
	}
	doJSONRequest(t, handler, http.MethodGet, "/customers/member-77/vouchers", nil, http.StatusOK, &vouchers)
	if len(vouchers) != 1 {
		t.Fatalf("expected one voucher, got %d", len(vouchers))
	}
}

func TestMetricsEndpointIsExposed(t *testing.T) {
	t.Parallel()

	server := newTestServer(t)
	handler := server.Handler()

	doJSONRequest(
		t,
		handler,
		http.MethodGet,
		"/healthz",
		nil,
		http.StatusOK,
		nil,
	)

	status, body := doRawRequest(t, handler, http.MethodGet, "/metrics")
	if status != http.StatusOK {
		t.Fatalf("expected 200 from metrics endpoint, got %d", status)
	}

	if !strings.Contains(body, "rich_domain_http_requests_total") {
		t.Fatalf("expected rich_domain_http_requests_total metric in output")
	}
}

func TestTraceIDContinuesFromRequestContext(t *testing.T) {
	t.Parallel()

	server := newTestServer(t)
	handler := server.Handler()

	const incomingTraceParent = "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
	status, traceID := doRawRequestWithTraceParent(t, handler, http.MethodGet, "/healthz", incomingTraceParent)

	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", status)
	}
	if traceID != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Fatalf("expected propagated trace id, got %q", traceID)
	}
}

func newTestServer(t *testing.T) *httpapi.Server {
	t.Helper()

	bus := eventbus.NewInMemoryEventBus()
	showRepository := inmemory.NewShowRepository()
	voucherRepository := inmemory.NewVoucherRepository()
	orderRepository := inmemory.NewOrderRepository()
	reservationRepository := inmemory.NewReservationRepository()

	reservationService := reservationapp.NewService(reservationRepository, bus)
	theaterService := theaterapp.NewService(showRepository, bus)
	kitchenService := kitchenapp.NewService(voucherRepository, orderRepository, bus)
	bus.Subscribe(theaterdomain.EventVIPSeatPurchased, kitchenService.HandleVIPSeatPurchased)

	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	t.Cleanup(func() {
		_ = tracerProvider.Shutdown(context.Background())
	})
	otel.SetTextMapPropagator(propagation.TraceContext{})

	server, err := httpapi.NewServer(
		reservationService,
		theaterService,
		kitchenService,
		zap.NewNop(),
		tracerProvider.Tracer("test"),
		prometheus.NewRegistry(),
	)
	if err != nil {
		t.Fatalf("new server: %v", err)
	}

	return server
}

func doJSONRequest(t *testing.T, handler http.Handler, method, path string, payload any, expectedStatus int, responseTarget any) {
	t.Helper()

	var body io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal payload: %v", err)
		}
		body = bytes.NewReader(encoded)
	}

	request, err := http.NewRequestWithContext(context.Background(), method, path, body)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	response := recorder.Result()
	t.Cleanup(func() { _ = response.Body.Close() })

	if response.StatusCode != expectedStatus {
		bodyBytes, _ := io.ReadAll(response.Body)
		t.Fatalf("expected status %d, got %d: %s", expectedStatus, response.StatusCode, string(bodyBytes))
	}

	if responseTarget == nil {
		return
	}

	if err := json.NewDecoder(response.Body).Decode(responseTarget); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}

func doRawRequest(t *testing.T, handler http.Handler, method, path string) (int, string) {
	t.Helper()

	request, err := http.NewRequestWithContext(context.Background(), method, path, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	return recorder.Code, recorder.Body.String()
}

func doRawRequestWithTraceParent(t *testing.T, handler http.Handler, method, path, traceParent string) (int, string) {
	t.Helper()

	request, err := http.NewRequestWithContext(context.Background(), method, path, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	request.Header.Set("traceparent", traceParent)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	return recorder.Code, recorder.Header().Get("X-Trace-ID")
}
