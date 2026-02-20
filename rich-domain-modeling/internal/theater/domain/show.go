package domain

import (
	"errors"
	"fmt"
	"time"

	shareddomain "richdomainmodeling/internal/shared/domain"
)

type SeatTier string

const (
	SeatTierStandard SeatTier = "standard"
	SeatTierVIP      SeatTier = "vip"
)

var (
	ErrInvalidShowID       = errors.New("show id is required")
	ErrInvalidShowTitle    = errors.New("show title is required")
	ErrNoSeatsConfigured   = errors.New("at least one seat is required")
	ErrDuplicateSeatNumber = errors.New("duplicate seat number")
	ErrInvalidSeatNumber   = errors.New("seat number is required")
	ErrInvalidSeatPrice    = errors.New("seat price must be positive")
	ErrInvalidCustomerID   = errors.New("customer id is required")
	ErrSeatNotFound        = errors.New("seat is not available")
)

type Seat struct {
	Number     string
	Tier       SeatTier
	PriceCents int
}

type Ticket struct {
	ID                 string
	ShowID             string
	CustomerID         string
	SeatNumber         string
	PriceCents         int
	IncludesFreeCoffee bool
	PurchasedAt        time.Time
}

type Show struct {
	recorder       shareddomain.EventRecorder
	id             string
	title          string
	startsAt       time.Time
	availableSeats map[string]Seat
	soldTickets    map[string]Ticket
}

func NewShow(id, title string, startsAt time.Time, seats []Seat) (*Show, error) {
	switch {
	case id == "":
		return nil, ErrInvalidShowID
	case title == "":
		return nil, ErrInvalidShowTitle
	case len(seats) == 0:
		return nil, ErrNoSeatsConfigured
	}

	availableSeats := make(map[string]Seat, len(seats))
	for _, seat := range seats {
		if seat.Number == "" {
			return nil, ErrInvalidSeatNumber
		}
		if seat.PriceCents <= 0 {
			return nil, ErrInvalidSeatPrice
		}
		if _, exists := availableSeats[seat.Number]; exists {
			return nil, fmt.Errorf("%w: %s", ErrDuplicateSeatNumber, seat.Number)
		}
		availableSeats[seat.Number] = seat
	}

	show := &Show{
		id:             id,
		title:          title,
		startsAt:       startsAt,
		availableSeats: availableSeats,
		soldTickets:    make(map[string]Ticket),
	}

	return show, nil
}

func (s *Show) ID() string {
	return s.id
}

func (s *Show) Title() string {
	return s.title
}

func (s *Show) StartsAt() time.Time {
	return s.startsAt
}

func (s *Show) RemainingSeats() int {
	return len(s.availableSeats)
}

func (s *Show) SoldTickets() []Ticket {
	tickets := make([]Ticket, 0, len(s.soldTickets))
	for _, ticket := range s.soldTickets {
		tickets = append(tickets, ticket)
	}
	return tickets
}

func (s *Show) PurchaseSeat(customerID, seatNumber string, purchasedAt time.Time) (Ticket, error) {
	if customerID == "" {
		return Ticket{}, ErrInvalidCustomerID
	}

	seat, exists := s.availableSeats[seatNumber]
	if !exists {
		return Ticket{}, fmt.Errorf("%w: %s", ErrSeatNotFound, seatNumber)
	}

	delete(s.availableSeats, seatNumber)

	ticketID := fmt.Sprintf("%s-%03d", s.id, len(s.soldTickets)+1)
	ticket := Ticket{
		ID:                 ticketID,
		ShowID:             s.id,
		CustomerID:         customerID,
		SeatNumber:         seat.Number,
		PriceCents:         seat.PriceCents,
		IncludesFreeCoffee: seat.Tier == SeatTierVIP,
		PurchasedAt:        purchasedAt,
	}

	s.soldTickets[ticketID] = ticket

	s.recorder.Record(NewSeatPurchased(ticketID, s.id, customerID, seat.Number, seat.PriceCents, purchasedAt))
	if seat.Tier == SeatTierVIP {
		s.recorder.Record(NewVIPSeatPurchased(ticketID, s.id, customerID, seat.Number, purchasedAt))
	}

	return ticket, nil
}

func (s *Show) PullEvents() []shareddomain.DomainEvent {
	return s.recorder.PullEvents()
}
