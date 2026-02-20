package domain

import (
	"errors"
	"fmt"
	"time"

	shareddomain "richdomainmodeling/internal/shared/domain"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusConfirmed Status = "confirmed"
	StatusCanceled  Status = "canceled"
)

var (
	ErrInvalidReservationID    = errors.New("reservation id is required")
	ErrInvalidWorkspaceID      = errors.New("workspace id is required")
	ErrInvalidMemberID         = errors.New("member id is required")
	ErrInvalidTimeRange        = errors.New("starts at must be before ends at")
	ErrReservationTimeConflict = errors.New("reservation time conflicts with an existing reservation")
	ErrCannotConfirm           = errors.New("only pending reservations can be confirmed")
	ErrCannotCancel            = errors.New("reservation is already canceled")
	ErrCancelReasonRequired    = errors.New("cancel reason is required")
)

type Reservation struct {
	recorder     shareddomain.EventRecorder
	id           string
	workspaceID  string
	memberID     string
	startsAt     time.Time
	endsAt       time.Time
	status       Status
	cancelReason string
}

func NewReservation(id, workspaceID, memberID string, startsAt, endsAt, now time.Time) (*Reservation, error) {
	switch {
	case id == "":
		return nil, ErrInvalidReservationID
	case workspaceID == "":
		return nil, ErrInvalidWorkspaceID
	case memberID == "":
		return nil, ErrInvalidMemberID
	case !startsAt.Before(endsAt):
		return nil, ErrInvalidTimeRange
	}

	reservation := &Reservation{
		id:          id,
		workspaceID: workspaceID,
		memberID:    memberID,
		startsAt:    startsAt,
		endsAt:      endsAt,
		status:      StatusPending,
	}

	reservation.recorder.Record(NewReservationCreated(id, workspaceID, memberID, startsAt, endsAt, now))

	return reservation, nil
}

func (r *Reservation) ID() string {
	return r.id
}

func (r *Reservation) WorkspaceID() string {
	return r.workspaceID
}

func (r *Reservation) MemberID() string {
	return r.memberID
}

func (r *Reservation) StartsAt() time.Time {
	return r.startsAt
}

func (r *Reservation) EndsAt() time.Time {
	return r.endsAt
}

func (r *Reservation) Status() Status {
	return r.status
}

func (r *Reservation) CancelReason() string {
	return r.cancelReason
}

func (r *Reservation) IsActive() bool {
	return r.status != StatusCanceled
}

func (r *Reservation) Overlaps(startsAt, endsAt time.Time) bool {
	return startsAt.Before(r.endsAt) && r.startsAt.Before(endsAt)
}

func (r *Reservation) Confirm(now time.Time) error {
	if r.status != StatusPending {
		return fmt.Errorf("%w: current status %q", ErrCannotConfirm, r.status)
	}

	r.status = StatusConfirmed
	r.recorder.Record(NewReservationConfirmed(r.id, now))

	return nil
}

func (r *Reservation) Cancel(reason string, now time.Time) error {
	if reason == "" {
		return ErrCancelReasonRequired
	}
	if r.status == StatusCanceled {
		return ErrCannotCancel
	}

	r.status = StatusCanceled
	r.cancelReason = reason
	r.recorder.Record(NewReservationCanceled(r.id, reason, now))

	return nil
}

func (r *Reservation) PullEvents() []shareddomain.DomainEvent {
	return r.recorder.PullEvents()
}
