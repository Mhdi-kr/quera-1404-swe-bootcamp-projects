package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"richdomainmodeling/internal/infrastructure/eventbus"
	"richdomainmodeling/internal/infrastructure/persistence/inmemory"
	"richdomainmodeling/internal/infrastructure/telemetry"
	"richdomainmodeling/internal/interfaces/httpapi"
	kitchenapp "richdomainmodeling/internal/kitchen/application"
	reservationapp "richdomainmodeling/internal/reservation/application"
	theaterapp "richdomainmodeling/internal/theater/application"
	theaterdomain "richdomainmodeling/internal/theater/domain"
)

const (
	defaultHTTPAddr    = ":8080"
	tracingServiceName = "rich-domain-modeling-httpapi"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	shutdownTracing, err := telemetry.SetupTracing(ctx, tracingServiceName)
	if err != nil {
		logger.Fatal("setup tracing", zap.Error(err))
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTracing(shutdownCtx); err != nil {
			logger.Error("shutdown tracing", zap.Error(err))
		}
	}()

	eventBus := eventbus.NewInMemoryEventBus()
	reservationRepository := inmemory.NewReservationRepository()
	showRepository := inmemory.NewShowRepository()
	voucherRepository := inmemory.NewVoucherRepository()
	orderRepository := inmemory.NewOrderRepository()

	reservationService := reservationapp.NewService(reservationRepository, eventBus)
	theaterService := theaterapp.NewService(showRepository, eventBus)
	kitchenService := kitchenapp.NewService(voucherRepository, orderRepository, eventBus)

	eventBus.Subscribe(theaterdomain.EventVIPSeatPurchased, kitchenService.HandleVIPSeatPurchased)

	registry := prometheus.NewRegistry()
	server, err := httpapi.NewServer(
		reservationService,
		theaterService,
		kitchenService,
		logger,
		otel.Tracer(tracingServiceName),
		registry,
	)
	if err != nil {
		logger.Fatal("setup http server", zap.Error(err))
	}

	addr := defaultHTTPAddr
	if envAddr := os.Getenv("HTTP_ADDR"); envAddr != "" {
		addr = envAddr
	}

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           server.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutdown http server", zap.Error(err))
		}
	}()

	logger.Info("http server started", zap.String("addr", addr))

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("http server stopped unexpectedly", zap.Error(err))
	}

	logger.Info("http server stopped")
}
