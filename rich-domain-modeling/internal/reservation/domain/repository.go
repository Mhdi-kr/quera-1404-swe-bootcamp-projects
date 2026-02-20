package domain

import (
	"context"
	"errors"
	"time"
)

var ErrReservationNotFound = errors.New("reservation not found")

type Repository interface {
	Save(context.Context, *Reservation) error
	GetByID(context.Context, string) (*Reservation, error)
	HasActiveReservationConflict(ctx context.Context, workspaceID string, startsAt, endsAt time.Time, excludingReservationID string) (bool, error)
}
