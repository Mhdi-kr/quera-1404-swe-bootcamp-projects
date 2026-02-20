package domain_test

import (
	"testing"
	"time"

	shareddomain "richdomainmodeling/internal/shared/domain"
	theaterdomain "richdomainmodeling/internal/theater/domain"
)

func TestShowPurchaseSeatPublishesVIPEventForVIPSeats(t *testing.T) {
	t.Parallel()

	show, err := theaterdomain.NewShow(
		"show-42",
		"Architecture Live",
		time.Now().UTC().Add(12*time.Hour),
		[]theaterdomain.Seat{{Number: "A2", Tier: theaterdomain.SeatTierVIP, PriceCents: 3200}},
	)
	if err != nil {
		t.Fatalf("new show: %v", err)
	}

	ticket, err := show.PurchaseSeat("user-1", "A2", time.Now().UTC())
	if err != nil {
		t.Fatalf("purchase seat: %v", err)
	}
	if !ticket.IncludesFreeCoffee {
		t.Fatalf("expected VIP ticket to include free coffee")
	}

	events := show.PullEvents()
	if len(events) != 2 {
		t.Fatalf("expected two events (purchase + vip purchase), got %d", len(events))
	}
	assertHasEvent(t, events, theaterdomain.EventVIPSeatPurchased)
}

func assertHasEvent(t *testing.T, events []shareddomain.DomainEvent, eventName string) {
	t.Helper()

	for _, event := range events {
		if event.EventName() == eventName {
			return
		}
	}

	t.Fatalf("event %q not found", eventName)
}
