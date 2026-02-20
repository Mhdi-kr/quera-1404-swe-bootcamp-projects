package domain

import "time"

const (
	EventSeatPurchased    = "theater.seat.purchased"
	EventVIPSeatPurchased = "theater.seat.vip.purchased"
)

type SeatPurchased struct {
	TicketID   string
	ShowID     string
	CustomerID string
	SeatNumber string
	PriceCents int
	OccurredOn time.Time
}

func NewSeatPurchased(ticketID, showID, customerID, seatNumber string, priceCents int, occurredOn time.Time) SeatPurchased {
	return SeatPurchased{
		TicketID:   ticketID,
		ShowID:     showID,
		CustomerID: customerID,
		SeatNumber: seatNumber,
		PriceCents: priceCents,
		OccurredOn: occurredOn,
	}
}

func (e SeatPurchased) EventName() string {
	return EventSeatPurchased
}

func (e SeatPurchased) OccurredAt() time.Time {
	return e.OccurredOn
}

type VIPSeatPurchased struct {
	TicketID   string
	ShowID     string
	CustomerID string
	SeatNumber string
	OccurredOn time.Time
}

func NewVIPSeatPurchased(ticketID, showID, customerID, seatNumber string, occurredOn time.Time) VIPSeatPurchased {
	return VIPSeatPurchased{
		TicketID:   ticketID,
		ShowID:     showID,
		CustomerID: customerID,
		SeatNumber: seatNumber,
		OccurredOn: occurredOn,
	}
}

func (e VIPSeatPurchased) EventName() string {
	return EventVIPSeatPurchased
}

func (e VIPSeatPurchased) OccurredAt() time.Time {
	return e.OccurredOn
}
