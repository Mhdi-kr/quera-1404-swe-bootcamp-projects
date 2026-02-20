package application

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	sharedapp "richdomainmodeling/internal/shared/application"
	shareddomain "richdomainmodeling/internal/shared/domain"
	theaterdomain "richdomainmodeling/internal/theater/domain"
)

type Service struct {
	repository theaterdomain.Repository
	eventBus   sharedapp.EventBus
	tracer     trace.Tracer
}

func NewService(repository theaterdomain.Repository, eventBus sharedapp.EventBus) *Service {
	return &Service{
		repository: repository,
		eventBus:   eventBus,
		tracer:     otel.Tracer("richdomainmodeling/internal/theater/application"),
	}
}

func (s *Service) ScheduleShow(ctx context.Context, showID, title string, startsAt time.Time, seats []theaterdomain.Seat) (*theaterdomain.Show, error) {
	_, domainSpan := s.tracer.Start(ctx, "domain.theater.show.new")
	domainSpan.SetAttributes(
		attribute.String("layer", "domain"),
		attribute.String("show.id", showID),
	)
	show, err := theaterdomain.NewShow(showID, title, startsAt, seats)
	if err != nil {
		domainSpan.RecordError(err)
		domainSpan.SetStatus(codes.Error, err.Error())
		domainSpan.End()
		return nil, err
	}
	domainSpan.End()

	if err := s.repository.Save(ctx, show); err != nil {
		return nil, err
	}

	return show, nil
}

func (s *Service) PurchaseSeat(ctx context.Context, showID, customerID, seatNumber string, purchasedAt time.Time) (theaterdomain.Ticket, error) {
	show, err := s.repository.GetByID(ctx, showID)
	if err != nil {
		return theaterdomain.Ticket{}, err
	}

	_, domainSpan := s.tracer.Start(ctx, "domain.theater.show.purchase_seat")
	domainSpan.SetAttributes(
		attribute.String("layer", "domain"),
		attribute.String("show.id", showID),
		attribute.String("seat.number", seatNumber),
		attribute.String("customer.id", customerID),
	)
	ticket, err := show.PurchaseSeat(customerID, seatNumber, purchasedAt)
	if err != nil {
		domainSpan.RecordError(err)
		domainSpan.SetStatus(codes.Error, err.Error())
		domainSpan.End()
		return theaterdomain.Ticket{}, err
	}
	domainSpan.End()

	if err := s.repository.Save(ctx, show); err != nil {
		return theaterdomain.Ticket{}, err
	}
	if err := s.publish(ctx, show.PullEvents()); err != nil {
		return theaterdomain.Ticket{}, err
	}

	return ticket, nil
}

func (s *Service) GetShow(ctx context.Context, showID string) (*theaterdomain.Show, error) {
	return s.repository.GetByID(ctx, showID)
}

func (s *Service) publish(ctx context.Context, events []shareddomain.DomainEvent) error {
	if s.eventBus == nil || len(events) == 0 {
		return nil
	}
	return s.eventBus.Publish(ctx, events...)
}
