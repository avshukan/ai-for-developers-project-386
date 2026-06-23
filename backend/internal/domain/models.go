// Package domain holds the Call Booking entities and the values that describe
// them. These types mirror the API contract in api/main.tsp (generated to
// openapi/openapi.yaml); the JSON tags are the wire field names from that
// contract. Business rules live in the store and slots packages, not here.
package domain

import "time"

// SlotDuration is the fixed length of a bookable slot in the MVP.
const SlotDuration = 30 * time.Minute

// SlotWindow is how far into the future guests can see and book slots.
const SlotWindow = 14 * 24 * time.Hour

// EventType is a bookable meeting type created by the host.
type EventType struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	DurationMinutes int    `json:"durationMinutes"`
}

// Availability is a time range published by the host. Slots are derived from
// it; it does not reserve time by itself.
type Availability struct {
	ID        string    `json:"id"`
	StartAt   time.Time `json:"startAt"`
	EndAt     time.Time `json:"endAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// SlotStatus is the availability state of a derived slot.
type SlotStatus string

const (
	// SlotAvailable marks a slot that can still be booked.
	SlotAvailable SlotStatus = "available"
	// SlotBooked marks a slot that is already taken. Kept for contract
	// fidelity; the list-slots endpoint only returns available slots.
	SlotBooked SlotStatus = "booked"
)

// Slot is a fixed 30-minute bookable interval derived from availability.
type Slot struct {
	StartAt time.Time  `json:"startAt"`
	EndAt   time.Time  `json:"endAt"`
	Status  SlotStatus `json:"status"`
}

// Booking is a confirmed reservation of a slot by a guest.
type Booking struct {
	ID          string    `json:"id"`
	EventTypeID string    `json:"eventTypeId"`
	StartAt     time.Time `json:"startAt"`
	EndAt       time.Time `json:"endAt"`
	GuestName   string    `json:"guestName"`
	GuestEmail  string    `json:"guestEmail"`
	CreatedAt   time.Time `json:"createdAt"`
}
