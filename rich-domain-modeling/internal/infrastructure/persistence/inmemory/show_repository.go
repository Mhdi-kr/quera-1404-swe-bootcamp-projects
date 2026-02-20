package inmemory

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	theaterdomain "richdomainmodeling/internal/theater/domain"
)

type ShowRepository struct {
	mu     sync.RWMutex
	shows  map[string]*theaterdomain.Show
	tracer trace.Tracer
}

func NewShowRepository() *ShowRepository {
	return &ShowRepository{
		shows:  make(map[string]*theaterdomain.Show),
		tracer: otel.Tracer("richdomainmodeling/internal/infrastructure/persistence/inmemory"),
	}
}

func (r *ShowRepository) Save(ctx context.Context, show *theaterdomain.Show) (err error) {
	_, span := r.tracer.Start(ctx, "repo.theater.show.save")
	span.SetAttributes(
		attribute.String("layer", "repo"),
		attribute.String("show.id", show.ID()),
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

	r.shows[show.ID()] = show
	return nil
}

func (r *ShowRepository) GetByID(ctx context.Context, showID string) (show *theaterdomain.Show, err error) {
	_, span := r.tracer.Start(ctx, "repo.theater.show.get_by_id")
	span.SetAttributes(
		attribute.String("layer", "repo"),
		attribute.String("show.id", showID),
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

	show, ok := r.shows[showID]
	if !ok {
		err = theaterdomain.ErrShowNotFound
		return nil, err
	}
	return show, nil
}
