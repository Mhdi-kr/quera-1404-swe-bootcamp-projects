package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"richdomainmodeling/internal/infrastructure/persistence/inmemory"
	reservationapp "richdomainmodeling/internal/reservation/application"
	reservationdomain "richdomainmodeling/internal/reservation/domain"
)

type spyRepository struct {
	conflictChecks int
	saves          int
}

func (r *spyRepository) Save(_ context.Context, _ *reservationdomain.Reservation) error {
	r.saves++
	return nil
}

func (r *spyRepository) GetByID(_ context.Context, _ string) (*reservationdomain.Reservation, error) {
	return nil, reservationdomain.ErrReservationNotFound
}

func (r *spyRepository) HasActiveReservationConflict(
	_ context.Context,
	_ string,
	_, _ time.Time,
	_ string,
) (bool, error) {
	r.conflictChecks++
	return false, nil
}

func TestReserveWorkspaceRejectsOverlappingReservations(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := inmemory.NewReservationRepository()
	service := reservationapp.NewService(repo, nil)
	base := time.Date(2026, 2, 20, 9, 0, 0, 0, time.UTC)

	_, err := service.ReserveWorkspace(
		ctx,
		"res-1",
		"workspace-1",
		"member-1",
		base,
		base.Add(2*time.Hour),
		base,
	)
	if err != nil {
		t.Fatalf("first reservation failed: %v", err)
	}

	_, err = service.ReserveWorkspace(
		ctx,
		"res-2",
		"workspace-1",
		"member-2",
		base.Add(time.Hour),
		base.Add(3*time.Hour),
		base.Add(time.Minute),
	)
	if !errors.Is(err, reservationdomain.ErrReservationTimeConflict) {
		t.Fatalf("expected overlap conflict error, got: %v", err)
	}
}

func TestReserveWorkspaceAllowsAdjacentTimeSlots(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := inmemory.NewReservationRepository()
	service := reservationapp.NewService(repo, nil)
	base := time.Date(2026, 2, 20, 9, 0, 0, 0, time.UTC)

	_, err := service.ReserveWorkspace(
		ctx,
		"res-1",
		"workspace-1",
		"member-1",
		base,
		base.Add(2*time.Hour),
		base,
	)
	if err != nil {
		t.Fatalf("first reservation failed: %v", err)
	}

	_, err = service.ReserveWorkspace(
		ctx,
		"res-2",
		"workspace-1",
		"member-2",
		base.Add(2*time.Hour),
		base.Add(4*time.Hour),
		base.Add(time.Minute),
	)
	if err != nil {
		t.Fatalf("adjacent reservation should be allowed, got: %v", err)
	}
}

func TestReserveWorkspaceAllowsSlotAfterCancellation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := inmemory.NewReservationRepository()
	service := reservationapp.NewService(repo, nil)
	base := time.Date(2026, 2, 20, 9, 0, 0, 0, time.UTC)

	reservation, err := service.ReserveWorkspace(
		ctx,
		"res-1",
		"workspace-1",
		"member-1",
		base,
		base.Add(2*time.Hour),
		base,
	)
	if err != nil {
		t.Fatalf("first reservation failed: %v", err)
	}

	if err := service.CancelReservation(ctx, reservation.ID(), "schedule changed", base.Add(30*time.Minute)); err != nil {
		t.Fatalf("cancel reservation failed: %v", err)
	}

	_, err = service.ReserveWorkspace(
		ctx,
		"res-2",
		"workspace-1",
		"member-2",
		base.Add(30*time.Minute),
		base.Add(90*time.Minute),
		base.Add(time.Hour),
	)
	if err != nil {
		t.Fatalf("overlap with canceled reservation should be allowed, got: %v", err)
	}
}

func TestReserveWorkspaceValidatesBeforeRepositoryCalls(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := &spyRepository{}
	service := reservationapp.NewService(repo, nil)
	base := time.Date(2026, 2, 20, 9, 0, 0, 0, time.UTC)

	_, err := service.ReserveWorkspace(
		ctx,
		"res-1",
		"workspace-1",
		"member-1",
		base,
		base,
		base,
	)
	if !errors.Is(err, reservationdomain.ErrInvalidTimeRange) {
		t.Fatalf("expected invalid time range error, got: %v", err)
	}
	if repo.conflictChecks != 0 {
		t.Fatalf("expected zero conflict checks, got: %d", repo.conflictChecks)
	}
	if repo.saves != 0 {
		t.Fatalf("expected zero saves, got: %d", repo.saves)
	}
}
