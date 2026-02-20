FROM golang:1.25.5 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/httpapi ./cmd/httpapi

FROM registry.docker.ir/alpine:3.20

WORKDIR /app
COPY --from=builder /out/httpapi /app/httpapi

EXPOSE 8080

ENV HTTP_ADDR=:8080
ENV TRACE_EXPORTER=zipkin
ENV ZIPKIN_ENDPOINT=http://tempo:9411/api/v2/spans

ENTRYPOINT ["/app/httpapi"]
