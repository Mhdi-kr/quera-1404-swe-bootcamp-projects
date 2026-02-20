package domain

import "time"

const (
	EventReservationCreated   = "reservation.created"
	EventReservationConfirmed = "reservation.confirmed"
	EventReservationCanceled  = "reservation.canceled"
)

type ReservationCreated struct {
	ReservationID string
	WorkspaceID   string
	MemberID      string
	StartsAt      time.Time
	EndsAt        time.Time
	OccurredOn    time.Time
}

func NewReservationCreated(reservationID, workspaceID, memberID string, startsAt, endsAt, occurredOn time.Time) ReservationCreated {
	return ReservationCreated{
		ReservationID: reservationID,
		WorkspaceID:   workspaceID,
		MemberID:      memberID,
		StartsAt:      startsAt,
		EndsAt:        endsAt,
		OccurredOn:    occurredOn,
	}
}

func (e ReservationCreated) EventName() string {
	return EventReservationCreated
}

func (e ReservationCreated) OccurredAt() time.Time {
	return e.OccurredOn
}

type ReservationConfirmed struct {
	ReservationID string
	OccurredOn    time.Time
}

func NewReservationConfirmed(reservationID string, occurredOn time.Time) ReservationConfirmed {
	return ReservationConfirmed{
		ReservationID: reservationID,
		OccurredOn:    occurredOn,
	}
}

func (e ReservationConfirmed) EventName() string {
	return EventReservationConfirmed
}

func (e ReservationConfirmed) OccurredAt() time.Time {
	return e.OccurredOn
}

type ReservationCanceled struct {
	ReservationID string
	Reason        string
	OccurredOn    time.Time
}

func NewReservationCanceled(reservationID, reason string, occurredOn time.Time) ReservationCanceled {
	return ReservationCanceled{
		ReservationID: reservationID,
		Reason:        reason,
		OccurredOn:    occurredOn,
	}
}

func (e ReservationCanceled) EventName() string {
	return EventReservationCanceled
}

func (e ReservationCanceled) OccurredAt() time.Time {
	return e.OccurredOn
}
