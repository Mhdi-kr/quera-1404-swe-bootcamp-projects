package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	kitchenapp "richdomainmodeling/internal/kitchen/application"
	kitchendomain "richdomainmodeling/internal/kitchen/domain"
	reservationapp "richdomainmodeling/internal/reservation/application"
	reservationdomain "richdomainmodeling/internal/reservation/domain"
	theaterapp "richdomainmodeling/internal/theater/application"
	theaterdomain "richdomainmodeling/internal/theater/domain"
)

const maxRequestBodyBytes = 1 << 20

type Server struct {
	reservationService *reservationapp.Service
	theaterService     *theaterapp.Service
	kitchenService     *kitchenapp.Service
	logger             *zap.Logger
	tracer             trace.Tracer
	registry           *prometheus.Registry
	requestTotal       *prometheus.CounterVec
	requestDuration    *prometheus.HistogramVec
}

func NewServer(
	reservationService *reservationapp.Service,
	theaterService *theaterapp.Service,
	kitchenService *kitchenapp.Service,
	logger *zap.Logger,
	tracer trace.Tracer,
	registry *prometheus.Registry,
) (*Server, error) {
	if logger == nil {
		logger = zap.NewNop()
	}
	if tracer == nil {
		tracer = otel.Tracer("richdomainmodeling/internal/interfaces/httpapi")
	}
	if registry == nil {
		registry = prometheus.NewRegistry()
	}

	requestTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rich_domain_http_requests_total",
			Help: "Total number of handled HTTP requests.",
		},
		[]string{"method", "route", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "rich_domain_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route", "status"},
	)

	requestTotal, err := registerCounterVec(registry, requestTotal)
	if err != nil {
		return nil, err
	}
	requestDuration, err = registerHistogramVec(registry, requestDuration)
	if err != nil {
		return nil, err
	}
	if err := registerCollector(registry, prometheus.NewGoCollector()); err != nil {
		return nil, err
	}
	if err := registerCollector(registry, prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{})); err != nil {
		return nil, err
	}

	return &Server{
		reservationService: reservationService,
		theaterService:     theaterService,
		kitchenService:     kitchenService,
		logger:             logger,
		tracer:             tracer,
		registry:           registry,
		requestTotal:       requestTotal,
		requestDuration:    requestDuration,
	}, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /healthz", s.wrap("GET /healthz", http.HandlerFunc(s.healthz)))
	mux.Handle("GET /metrics", s.wrap("GET /metrics", promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{})))

	mux.Handle("POST /reservations", s.wrap("POST /reservations", http.HandlerFunc(s.createReservation)))
	mux.Handle("GET /reservations/{reservation_id}", s.wrap("GET /reservations/{reservation_id}", http.HandlerFunc(s.getReservation)))
	mux.Handle("POST /reservations/{reservation_id}/confirm", s.wrap("POST /reservations/{reservation_id}/confirm", http.HandlerFunc(s.confirmReservation)))
	mux.Handle("POST /reservations/{reservation_id}/cancel", s.wrap("POST /reservations/{reservation_id}/cancel", http.HandlerFunc(s.cancelReservation)))

	mux.Handle("POST /shows", s.wrap("POST /shows", http.HandlerFunc(s.scheduleShow)))
	mux.Handle("GET /shows/{show_id}", s.wrap("GET /shows/{show_id}", http.HandlerFunc(s.getShow)))
	mux.Handle("POST /shows/{show_id}/purchase", s.wrap("POST /shows/{show_id}/purchase", http.HandlerFunc(s.purchaseSeat)))

	mux.Handle("GET /customers/{customer_id}/vouchers", s.wrap("GET /customers/{customer_id}/vouchers", http.HandlerFunc(s.listVouchersByCustomer)))
	mux.Handle("POST /kitchen/orders/paid", s.wrap("POST /kitchen/orders/paid", http.HandlerFunc(s.placePaidOrder)))
	mux.Handle("POST /kitchen/vouchers/{voucher_id}/redeem", s.wrap("POST /kitchen/vouchers/{voucher_id}/redeem", http.HandlerFunc(s.redeemVoucher)))

	return mux
}

func (s *Server) wrap(route string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		writer := &statusWriter{ResponseWriter: w}
		// Continue an upstream trace if a valid traceparent is provided.
		parentCtx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		ctx, span := s.tracer.Start(parentCtx, route, trace.WithSpanKind(trace.SpanKindServer))
		spanCtx := span.SpanContext()
		if spanCtx.IsValid() {
			writer.Header().Set("X-Trace-ID", spanCtx.TraceID().String())
		}
		r = r.WithContext(ctx)

		defer func() {
			if panicValue := recover(); panicValue != nil {
				span.RecordError(fmt.Errorf("panic: %v", panicValue))
				span.SetStatus(codes.Error, "panic recovered")

				s.logger.Error(
					"http panic recovered",
					zap.Any("panic", panicValue),
					zap.ByteString("stack", debug.Stack()),
					zap.String("route", route),
					zap.String("path", r.URL.Path),
				)

				if writer.status == 0 {
					_ = writeJSON(writer, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
				} else if writer.status < http.StatusInternalServerError {
					writer.status = http.StatusInternalServerError
				}
			}

			if writer.status == 0 {
				writer.status = http.StatusOK
			}

			duration := time.Since(startedAt)
			statusCode := strconv.Itoa(writer.status)

			s.requestTotal.WithLabelValues(r.Method, route, statusCode).Inc()
			s.requestDuration.WithLabelValues(r.Method, route, statusCode).Observe(duration.Seconds())

			span.SetAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.route", route),
				attribute.String("url.path", r.URL.Path),
				attribute.Int("http.status_code", writer.status),
			)
			if writer.status >= http.StatusInternalServerError {
				span.SetStatus(codes.Error, http.StatusText(writer.status))
			} else {
				span.SetStatus(codes.Ok, "")
			}

			logFields := []zap.Field{
				zap.String("method", r.Method),
				zap.String("route", route),
				zap.String("path", r.URL.Path),
				zap.Int("status", writer.status),
				zap.Duration("duration", duration),
			}
			if traceID := spanCtx.TraceID().String(); traceID != "" && traceID != "00000000000000000000000000000000" {
				logFields = append(logFields, zap.String("trace_id", traceID))
			}

			if writer.status >= http.StatusInternalServerError {
				s.logger.Error("http request completed", logFields...)
			} else {
				s.logger.Info("http request completed", logFields...)
			}

			span.End()
		}()

		next.ServeHTTP(writer, r)
	})
}

func (s *Server) healthz(w http.ResponseWriter, _ *http.Request) {
	_ = writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) createReservation(w http.ResponseWriter, r *http.Request) {
	var req createReservationRequest
	if err := decodeJSON(w, r, &req); err != nil {
		s.writeBadRequest(w, err)
		return
	}

	reservation, err := s.reservationService.ReserveWorkspace(
		r.Context(),
		req.ReservationID,
		req.WorkspaceID,
		req.MemberID,
		req.StartsAt.UTC(),
		req.EndsAt.UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		s.writeDomainError(w, err)
		return
	}

	_ = writeJSON(w, http.StatusCreated, toReservationResponse(reservation))
}

func (s *Server) getReservation(w http.ResponseWriter, r *http.Request) {
	reservationID := r.PathValue("reservation_id")
	reservation, err := s.reservationService.GetReservation(r.Context(), reservationID)
	if err != nil {
		s.writeDomainError(w, err)
		return
	}

	_ = writeJSON(w, http.StatusOK, toReservationResponse(reservation))
}

func (s *Server) confirmReservation(w http.ResponseWriter, r *http.Request) {
	reservationID := r.PathValue("reservation_id")
	if err := s.reservationService.ConfirmReservation(r.Context(), reservationID, time.Now().UTC()); err != nil {
		s.writeDomainError(w, err)
		return
	}

	reservation, err := s.reservationService.GetReservation(r.Context(), reservationID)
	if err != nil {
		s.writeDomainError(w, err)
		return
	}

	_ = writeJSON(w, http.StatusOK, toReservationResponse(reservation))
}

func (s *Server) cancelReservation(w http.ResponseWriter, r *http.Request) {
	reservationID := r.PathValue("reservation_id")

	var req cancelReservationRequest
	if err := decodeJSON(w, r, &req); err != nil {
		s.writeBadRequest(w, err)
		return
	}

	if err := s.reservationService.CancelReservation(r.Context(), reservationID, req.Reason, time.Now().UTC()); err != nil {
		s.writeDomainError(w, err)
		return
	}

	reservation, err := s.reservationService.GetReservation(r.Context(), reservationID)
	if err != nil {
		s.writeDomainError(w, err)
		return
	}

	_ = writeJSON(w, http.StatusOK, toReservationResponse(reservation))
}

func (s *Server) scheduleShow(w http.ResponseWriter, r *http.Request) {
	var req scheduleShowRequest
	if err := decodeJSON(w, r, &req); err != nil {
		s.writeBadRequest(w, err)
		return
	}

	seats := make([]theaterdomain.Seat, 0, len(req.Seats))
	for _, seat := range req.Seats {
		tier, err := parseSeatTier(seat.Tier)
		if err != nil {
			s.writeBadRequest(w, err)
			return
		}

		seats = append(seats, theaterdomain.Seat{
			Number:     seat.Number,
			Tier:       tier,
			PriceCents: seat.PriceCents,
		})
	}

	show, err := s.theaterService.ScheduleShow(
		r.Context(),
		req.ShowID,
		req.Title,
		req.StartsAt.UTC(),
		seats,
	)
	if err != nil {
		s.writeDomainError(w, err)
		return
	}

	_ = writeJSON(w, http.StatusCreated, toShowResponse(show))
}

func (s *Server) getShow(w http.ResponseWriter, r *http.Request) {
	showID := r.PathValue("show_id")
	show, err := s.theaterService.GetShow(r.Context(), showID)
	if err != nil {
		s.writeDomainError(w, err)
		return
	}

	_ = writeJSON(w, http.StatusOK, toShowResponse(show))
}

func (s *Server) purchaseSeat(w http.ResponseWriter, r *http.Request) {
	showID := r.PathValue("show_id")

	var req purchaseSeatRequest
	if err := decodeJSON(w, r, &req); err != nil {
		s.writeBadRequest(w, err)
		return
	}

	ticket, err := s.theaterService.PurchaseSeat(
		r.Context(),
		showID,
		req.CustomerID,
		req.SeatNumber,
		time.Now().UTC(),
	)
	if err != nil {
		s.writeDomainError(w, err)
		return
	}

	_ = writeJSON(w, http.StatusCreated, toTicketResponse(ticket))
}

func (s *Server) listVouchersByCustomer(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customer_id")
	vouchers, err := s.kitchenService.ListVouchersByCustomer(r.Context(), customerID)
	if err != nil {
		s.writeDomainError(w, err)
		return
	}

	response := make([]voucherResponse, 0, len(vouchers))
	for _, voucher := range vouchers {
		response = append(response, toVoucherResponse(voucher))
	}

	_ = writeJSON(w, http.StatusOK, response)
}

func (s *Server) placePaidOrder(w http.ResponseWriter, r *http.Request) {
	var req placePaidOrderRequest
	if err := decodeJSON(w, r, &req); err != nil {
		s.writeBadRequest(w, err)
		return
	}

	order, err := s.kitchenService.PlacePaidOrder(
		r.Context(),
		req.OrderID,
		req.CustomerID,
		req.Drink,
		req.PriceCents,
		time.Now().UTC(),
	)
	if err != nil {
		s.writeDomainError(w, err)
		return
	}

	_ = writeJSON(w, http.StatusCreated, toOrderResponse(order))
}

func (s *Server) redeemVoucher(w http.ResponseWriter, r *http.Request) {
	voucherID := r.PathValue("voucher_id")

	var req redeemVoucherRequest
	if err := decodeJSON(w, r, &req); err != nil {
		s.writeBadRequest(w, err)
		return
	}

	order, err := s.kitchenService.RedeemVoucher(
		r.Context(),
		voucherID,
		req.OrderID,
		req.Drink,
		time.Now().UTC(),
	)
	if err != nil {
		s.writeDomainError(w, err)
		return
	}

	_ = writeJSON(w, http.StatusCreated, toOrderResponse(order))
}

func (s *Server) writeBadRequest(w http.ResponseWriter, err error) {
	_ = writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
}

func (s *Server) writeDomainError(w http.ResponseWriter, err error) {
	status := mapDomainErrorToStatus(err)
	message := err.Error()
	if status == http.StatusInternalServerError {
		message = "internal server error"
	}

	_ = writeJSON(w, status, errorResponse{Error: message})
}

func mapDomainErrorToStatus(err error) int {
	switch {
	case errors.Is(err, reservationdomain.ErrReservationNotFound),
		errors.Is(err, theaterdomain.ErrShowNotFound),
		errors.Is(err, kitchendomain.ErrVoucherNotFound),
		errors.Is(err, kitchendomain.ErrOrderNotFound):
		return http.StatusNotFound

	case errors.Is(err, reservationdomain.ErrReservationTimeConflict),
		errors.Is(err, kitchendomain.ErrVoucherAlreadyExists),
		errors.Is(err, kitchendomain.ErrOrderAlreadyExists):
		return http.StatusConflict

	case errors.Is(err, reservationdomain.ErrCannotConfirm),
		errors.Is(err, reservationdomain.ErrCannotCancel),
		errors.Is(err, theaterdomain.ErrSeatNotFound),
		errors.Is(err, kitchendomain.ErrVoucherNotRedeemable):
		return http.StatusUnprocessableEntity

	case errors.Is(err, reservationdomain.ErrInvalidReservationID),
		errors.Is(err, reservationdomain.ErrInvalidWorkspaceID),
		errors.Is(err, reservationdomain.ErrInvalidMemberID),
		errors.Is(err, reservationdomain.ErrInvalidTimeRange),
		errors.Is(err, reservationdomain.ErrCancelReasonRequired),
		errors.Is(err, theaterdomain.ErrInvalidShowID),
		errors.Is(err, theaterdomain.ErrInvalidShowTitle),
		errors.Is(err, theaterdomain.ErrNoSeatsConfigured),
		errors.Is(err, theaterdomain.ErrDuplicateSeatNumber),
		errors.Is(err, theaterdomain.ErrInvalidSeatNumber),
		errors.Is(err, theaterdomain.ErrInvalidSeatPrice),
		errors.Is(err, theaterdomain.ErrInvalidCustomerID),
		errors.Is(err, kitchendomain.ErrInvalidVoucherID),
		errors.Is(err, kitchendomain.ErrInvalidVoucherSource),
		errors.Is(err, kitchendomain.ErrInvalidVoucherOwner),
		errors.Is(err, kitchendomain.ErrInvalidOrderID),
		errors.Is(err, kitchendomain.ErrInvalidDrink),
		errors.Is(err, kitchendomain.ErrInvalidOrderPrice),
		errors.Is(err, kitchendomain.ErrInvalidOrderOwner):
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}

func parseSeatTier(tier string) (theaterdomain.SeatTier, error) {
	switch tier {
	case string(theaterdomain.SeatTierStandard):
		return theaterdomain.SeatTierStandard, nil
	case string(theaterdomain.SeatTierVIP):
		return theaterdomain.SeatTierVIP, nil
	default:
		return "", fmt.Errorf("invalid seat tier: %s", tier)
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("request body must contain a single JSON object")
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	return encoder.Encode(payload)
}

func registerCounterVec(registry *prometheus.Registry, collector *prometheus.CounterVec) (*prometheus.CounterVec, error) {
	if err := registry.Register(collector); err != nil {
		alreadyRegistered, ok := err.(prometheus.AlreadyRegisteredError)
		if !ok {
			return nil, err
		}

		existing, ok := alreadyRegistered.ExistingCollector.(*prometheus.CounterVec)
		if !ok {
			return nil, err
		}
		return existing, nil
	}

	return collector, nil
}

func registerHistogramVec(registry *prometheus.Registry, collector *prometheus.HistogramVec) (*prometheus.HistogramVec, error) {
	if err := registry.Register(collector); err != nil {
		alreadyRegistered, ok := err.(prometheus.AlreadyRegisteredError)
		if !ok {
			return nil, err
		}

		existing, ok := alreadyRegistered.ExistingCollector.(*prometheus.HistogramVec)
		if !ok {
			return nil, err
		}
		return existing, nil
	}

	return collector, nil
}

func registerCollector(registry *prometheus.Registry, collector prometheus.Collector) error {
	if err := registry.Register(collector); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
			return nil
		}
		return err
	}
	return nil
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(payload []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(payload)
}

type errorResponse struct {
	Error string `json:"error"`
}

type createReservationRequest struct {
	ReservationID string    `json:"reservation_id"`
	WorkspaceID   string    `json:"workspace_id"`
	MemberID      string    `json:"member_id"`
	StartsAt      time.Time `json:"starts_at"`
	EndsAt        time.Time `json:"ends_at"`
}

type cancelReservationRequest struct {
	Reason string `json:"reason"`
}

type reservationResponse struct {
	ID           string                   `json:"id"`
	WorkspaceID  string                   `json:"workspace_id"`
	MemberID     string                   `json:"member_id"`
	StartsAt     time.Time                `json:"starts_at"`
	EndsAt       time.Time                `json:"ends_at"`
	Status       reservationdomain.Status `json:"status"`
	CancelReason string                   `json:"cancel_reason,omitempty"`
}

func toReservationResponse(reservation *reservationdomain.Reservation) reservationResponse {
	return reservationResponse{
		ID:           reservation.ID(),
		WorkspaceID:  reservation.WorkspaceID(),
		MemberID:     reservation.MemberID(),
		StartsAt:     reservation.StartsAt(),
		EndsAt:       reservation.EndsAt(),
		Status:       reservation.Status(),
		CancelReason: reservation.CancelReason(),
	}
}

type scheduleShowRequest struct {
	ShowID   string             `json:"show_id"`
	Title    string             `json:"title"`
	StartsAt time.Time          `json:"starts_at"`
	Seats    []scheduleShowSeat `json:"seats"`
}

type scheduleShowSeat struct {
	Number     string `json:"number"`
	Tier       string `json:"tier"`
	PriceCents int    `json:"price_cents"`
}

type showResponse struct {
	ID             string           `json:"id"`
	Title          string           `json:"title"`
	StartsAt       time.Time        `json:"starts_at"`
	RemainingSeats int              `json:"remaining_seats"`
	SoldTickets    []ticketResponse `json:"sold_tickets"`
}

func toShowResponse(show *theaterdomain.Show) showResponse {
	soldTickets := show.SoldTickets()
	responseTickets := make([]ticketResponse, 0, len(soldTickets))
	for _, ticket := range soldTickets {
		responseTickets = append(responseTickets, toTicketResponse(ticket))
	}

	return showResponse{
		ID:             show.ID(),
		Title:          show.Title(),
		StartsAt:       show.StartsAt(),
		RemainingSeats: show.RemainingSeats(),
		SoldTickets:    responseTickets,
	}
}

type purchaseSeatRequest struct {
	CustomerID string `json:"customer_id"`
	SeatNumber string `json:"seat_number"`
}

type ticketResponse struct {
	ID                 string    `json:"id"`
	ShowID             string    `json:"show_id"`
	CustomerID         string    `json:"customer_id"`
	SeatNumber         string    `json:"seat_number"`
	PriceCents         int       `json:"price_cents"`
	IncludesFreeCoffee bool      `json:"includes_free_coffee"`
	PurchasedAt        time.Time `json:"purchased_at"`
}

func toTicketResponse(ticket theaterdomain.Ticket) ticketResponse {
	return ticketResponse{
		ID:                 ticket.ID,
		ShowID:             ticket.ShowID,
		CustomerID:         ticket.CustomerID,
		SeatNumber:         ticket.SeatNumber,
		PriceCents:         ticket.PriceCents,
		IncludesFreeCoffee: ticket.IncludesFreeCoffee,
		PurchasedAt:        ticket.PurchasedAt,
	}
}

type placePaidOrderRequest struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	Drink      string `json:"drink"`
	PriceCents int    `json:"price_cents"`
}

type redeemVoucherRequest struct {
	OrderID string `json:"order_id"`
	Drink   string `json:"drink"`
}

type orderResponse struct {
	ID            string    `json:"id"`
	CustomerID    string    `json:"customer_id"`
	Drink         string    `json:"drink"`
	PriceCents    int       `json:"price_cents"`
	Complimentary bool      `json:"complimentary"`
	CreatedAt     time.Time `json:"created_at"`
}

func toOrderResponse(order *kitchendomain.CoffeeOrder) orderResponse {
	return orderResponse{
		ID:            order.ID(),
		CustomerID:    order.CustomerID(),
		Drink:         order.Drink(),
		PriceCents:    order.PriceCents(),
		Complimentary: order.Complimentary(),
		CreatedAt:     order.CreatedAt(),
	}
}

type voucherResponse struct {
	ID         string                      `json:"id"`
	CustomerID string                      `json:"customer_id"`
	Source     string                      `json:"source"`
	Status     kitchendomain.VoucherStatus `json:"status"`
	IssuedAt   time.Time                   `json:"issued_at"`
	RedeemedAt *time.Time                  `json:"redeemed_at,omitempty"`
}

func toVoucherResponse(voucher *kitchendomain.CoffeeVoucher) voucherResponse {
	return voucherResponse{
		ID:         voucher.ID(),
		CustomerID: voucher.CustomerID(),
		Source:     voucher.Source(),
		Status:     voucher.Status(),
		IssuedAt:   voucher.IssuedAt(),
		RedeemedAt: voucher.RedeemedAt(),
	}
}
