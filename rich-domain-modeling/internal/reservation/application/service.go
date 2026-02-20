package application

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	reservationdomain "richdomainmodeling/internal/reservation/domain"
	sharedapp "richdomainmodeling/internal/shared/application"
	shareddomain "richdomainmodeling/internal/shared/domain"
)

type Service struct {
	repository reservationdomain.Repository
	eventBus   sharedapp.EventBus
	tracer     trace.Tracer
}

func NewService(repository reservationdomain.Repository, eventBus sharedapp.EventBus) *Service {
	return &Service{
		repository: repository,
		eventBus:   eventBus,
		tracer:     otel.Tracer("richdomainmodeling/internal/reservation/application"),
	}
}

func (s *Service) ReserveWorkspace(
	ctx context.Context,
	reservationID string,
	workspaceID string,
	memberID string,
	startsAt time.Time,
	endsAt time.Time,
	now time.Time,
) (*reservationdomain.Reservation, error) {
	_, domainSpan := s.tracer.Start(ctx, "domain.reservation.new")
	domainSpan.SetAttributes(
		attribute.String("layer", "domain"),
		attribute.String("reservation.id", reservationID),
		attribute.String("reservation.workspace_id", workspaceID),
	)

	reservation, err := reservationdomain.NewReservation(reservationID, workspaceID, memberID, startsAt, endsAt, now)
	if err != nil {
		domainSpan.RecordError(err)
		domainSpan.SetStatus(codes.Error, err.Error())
		domainSpan.End()
		return nil, err
	}
	domainSpan.End()

	hasConflict, err := s.repository.HasActiveReservationConflict(
		ctx,
		reservation.WorkspaceID(),
		reservation.StartsAt(),
		reservation.EndsAt(),
		reservation.ID(),
	)
	if err != nil {
		return nil, err
	}
	if hasConflict {
		return nil, reservationdomain.ErrReservationTimeConflict
	}

	if err := s.repository.Save(ctx, reservation); err != nil {
		return nil, err
	}

	if err := s.publish(ctx, reservation.PullEvents()); err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *Service) ConfirmReservation(ctx context.Context, reservationID string, now time.Time) error {
	reservation, err := s.repository.GetByID(ctx, reservationID)
	if err != nil {
		return err
	}

	_, domainSpan := s.tracer.Start(ctx, "domain.reservation.confirm")
	domainSpan.SetAttributes(
		attribute.String("layer", "domain"),
		attribute.String("reservation.id", reservationID),
	)
	if err := reservation.Confirm(now); err != nil {
		domainSpan.RecordError(err)
		domainSpan.SetStatus(codes.Error, err.Error())
		domainSpan.End()
		return err
	}
	domainSpan.End()

	if err := s.repository.Save(ctx, reservation); err != nil {
		return err
	}

	return s.publish(ctx, reservation.PullEvents())
}

func (s *Service) CancelReservation(ctx context.Context, reservationID, reason string, now time.Time) error {
	reservation, err := s.repository.GetByID(ctx, reservationID)
	if err != nil {
		return err
	}

	_, domainSpan := s.tracer.Start(ctx, "domain.reservation.cancel")
	domainSpan.SetAttributes(
		attribute.String("layer", "domain"),
		attribute.String("reservation.id", reservationID),
	)
	if err := reservation.Cancel(reason, now); err != nil {
		domainSpan.RecordError(err)
		domainSpan.SetStatus(codes.Error, err.Error())
		domainSpan.End()
		return err
	}
	domainSpan.End()

	if err := s.repository.Save(ctx, reservation); err != nil {
		return err
	}

	return s.publish(ctx, reservation.PullEvents())
}

func (s *Service) GetReservation(ctx context.Context, reservationID string) (*reservationdomain.Reservation, error) {
	return s.repository.GetByID(ctx, reservationID)
}

func (s *Service) publish(ctx context.Context, events []shareddomain.DomainEvent) error {
	if s.eventBus == nil || len(events) == 0 {
		return nil
	}
	return s.eventBus.Publish(ctx, events...)
}
