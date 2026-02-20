# Rich Domain Modeling in Go (DDD)

This project demonstrates DDD structure with three bounded contexts:

- `reservation` domain for shared workspace reservations.
- `theater` domain for show seat sales.
- `kitchen` subdomain for coffee vouchers and coffee orders.

When a VIP seat is purchased in the theater domain, a domain event is published and handled by the kitchen application service to issue a complimentary coffee voucher.

## Project Layout

- `cmd/app`: Composition root and end-to-end demo.
- `cmd/httpapi`: Thin HTTP transport with observability (Zap + OTel + Prometheus).
- `internal/shared`: Shared kernel (`DomainEvent`, event bus contracts).
- `internal/reservation`: Reservation bounded context.
- `internal/theater`: Theater bounded context.
- `internal/kitchen`: Kitchen bounded context.
- `internal/infrastructure`: In-memory event bus and repositories.
- `internal/integration`: Integration test for cross-domain event flow.

## Run Demo

```bash
go test ./...
go run ./cmd/app
```

## Run HTTP API

```bash
go run ./cmd/httpapi
```

Optional env var:

- `HTTP_ADDR` (default `:8080`)
- `TRACE_EXPORTER` (`stdout` or `zipkin`, default `stdout`)
- `ZIPKIN_ENDPOINT` (used when `TRACE_EXPORTER=zipkin`)

Observability:

- Logs: structured JSON logs via Zap.
- Traces: OpenTelemetry spans exported to stdout (or Zipkin-compatible endpoint).
- Metrics: Prometheus endpoint at `GET /metrics`.

## Run Full Observability Stack

This repository includes Docker Compose for:

- Grafana (`http://localhost:3000`, `admin` / `admin`)
- Prometheus (`http://localhost:9090`)
- Loki (`http://localhost:3100`)
- Tempo (`http://localhost:3200`)

Run:

```bash
docker compose up --build
```

Notes:

- The app is configured in Compose with `TRACE_EXPORTER=zipkin` and sends traces to Tempo.
- Prometheus scrapes app metrics from `/metrics`.
- Promtail ships container logs to Loki.

## HTTP Endpoints

- `GET /healthz`
- `POST /reservations`
- `GET /reservations/{reservation_id}`
- `POST /reservations/{reservation_id}/confirm`
- `POST /reservations/{reservation_id}/cancel`
- `POST /shows`
- `GET /shows/{show_id}`
- `POST /shows/{show_id}/purchase`
- `GET /customers/{customer_id}/vouchers`
- `POST /kitchen/orders/paid`
- `POST /kitchen/vouchers/{voucher_id}/redeem`
