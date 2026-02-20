package inmemory

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	reservationdomain "richdomainmodeling/internal/reservation/domain"
)

type ReservationRepository struct {
	mu           sync.RWMutex
	reservations map[string]*reservationdomain.Reservation
	tracer       trace.Tracer
}

func NewReservationRepository() *ReservationRepository {
	return &ReservationRepository{
		reservations: make(map[string]*reservationdomain.Reservation),
		tracer:       otel.Tracer("richdomainmodeling/internal/infrastructure/persistence/inmemory"),
	}
}

func (r *ReservationRepository) Save(ctx context.Context, reservation *reservationdomain.Reservation) (err error) {
	_, span := r.tracer.Start(ctx, "repo.reservation.save")
	span.SetAttributes(
		attribute.String("layer", "repo"),
		attribute.String("reservation.id", reservation.ID()),
	)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	r.mu.Lock()
	defer r.mu.Unlock()

	r.reservations[reservation.ID()] = reservation
	return nil
}

func (r *ReservationRepository) GetByID(ctx context.Context, reservationID string) (reservation *reservationdomain.Reservation, err error) {
	_, span := r.tracer.Start(ctx, "repo.reservation.get_by_id")
	span.SetAttributes(
		attribute.String("layer", "repo"),
		attribute.String("reservation.id", reservationID),
	)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	r.mu.RLock()
	defer r.mu.RUnlock()

	reservation, ok := r.reservations[reservationID]
	if !ok {
		err = reservationdomain.ErrReservationNotFound
		return nil, err
	}
	return reservation, nil
}

func (r *ReservationRepository) HasActiveReservationConflict(
	ctx context.Context,
	workspaceID string,
	startsAt, endsAt time.Time,
	excludingReservationID string,
) (hasConflict bool, err error) {
	_, span := r.tracer.Start(ctx, "repo.reservation.has_active_conflict")
	span.SetAttributes(
		attribute.String("layer", "repo"),
		attribute.String("reservation.workspace_id", workspaceID),
	)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.SetAttributes(attribute.Bool("reservation.has_conflict", hasConflict))
		span.End()
	}()

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, reservation := range r.reservations {
		if reservation.ID() == excludingReservationID {
			continue
		}
		if reservation.WorkspaceID() != workspaceID {
			continue
		}
		if reservation.IsActive() && reservation.Overlaps(startsAt, endsAt) {
			return true, nil
		}
	}

	return false, nil
}
