// Package store is the in-memory data store for the Call Booking MVP.
//
// Storage choice is recorded in docs/adr/0002-backend-stack-and-storage.md:
// an in-memory store is sufficient because the MVP allows data to reset on
// restart. The store is the single owner of the no-double-booking invariant —
// CreateBooking checks for conflicts and inserts under one lock, so concurrent
// requests for the same time cannot both succeed.
package store

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/avshukan/ai-for-developers-project-386/backend/internal/domain"
	"github.com/avshukan/ai-for-developers-project-386/backend/internal/slots"
)

// Errors returned by CreateBooking, mapped to HTTP status codes by the HTTP
// layer (404 / 400 / 409 respectively).
var (
	// ErrEventTypeNotFound means the referenced event type does not exist.
	ErrEventTypeNotFound = errors.New("event type not found")
	// ErrSlotUnavailable means the requested start time is not a bookable
	// slot (in the past, outside availability, or misaligned).
	ErrSlotUnavailable = errors.New("selected slot is not available")
	// ErrBookingConflict means the slot is valid but already taken by an
	// overlapping booking.
	ErrBookingConflict = errors.New("booking conflict")
)

// Store holds event types, availability ranges and bookings in memory.
// The zero value is not usable; construct one with New.
type Store struct {
	mu sync.Mutex

	eventTypes    map[string]domain.EventType
	availabilites []domain.Availability
	bookings      []domain.Booking

	eventSeq        int
	availabilitySeq int
	bookingSeq      int

	// now is the clock, injectable so tests can pin server time.
	now func() time.Time
}

// New returns an empty store that uses the real wall clock.
func New() *Store {
	return NewWithClock(time.Now)
}

// NewWithClock returns an empty store driven by the given clock. Used by tests
// to make slot and "past" rules deterministic.
func NewWithClock(now func() time.Time) *Store {
	return &Store{
		eventTypes: make(map[string]domain.EventType),
		now:        now,
	}
}

// CreateEventTypeInput is the validated input for CreateEventType.
type CreateEventTypeInput struct {
	Title           string
	Description     string
	DurationMinutes int
}

// CreateEventType stores a new event type and returns it with a generated id.
func (s *Store) CreateEventType(in CreateEventTypeInput) domain.EventType {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.eventSeq++
	et := domain.EventType{
		ID:              "evt_" + strconv.Itoa(s.eventSeq),
		Title:           strings.TrimSpace(in.Title),
		Description:     strings.TrimSpace(in.Description),
		DurationMinutes: in.DurationMinutes,
	}
	s.eventTypes[et.ID] = et
	return et
}

// ListEventTypes returns all event types, ordered by creation (id) for a
// stable response.
func (s *Store) ListEventTypes() []domain.EventType {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]domain.EventType, 0, len(s.eventTypes))
	for _, et := range s.eventTypes {
		out = append(out, et)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// EventTypeExists reports whether an event type with the given id is stored.
func (s *Store) EventTypeExists(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.eventTypes[id]
	return ok
}

// CreateAvailabilityInput is the validated input for CreateAvailability.
type CreateAvailabilityInput struct {
	StartAt time.Time
	EndAt   time.Time
}

// CreateAvailability stores a new availability range and returns it with a
// generated id. Times are normalized to whole-second UTC.
func (s *Store) CreateAvailability(in CreateAvailabilityInput) domain.Availability {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.availabilitySeq++
	a := domain.Availability{
		ID:        "avl_" + strconv.Itoa(s.availabilitySeq),
		StartAt:   normalize(in.StartAt),
		EndAt:     normalize(in.EndAt),
		CreatedAt: normalize(s.now()),
	}
	s.availabilites = append(s.availabilites, a)
	return a
}

// AvailableSlots returns the available slots for the next 14 days, derived from
// stored availability minus existing bookings.
func (s *Store) AvailableSlots() []domain.Slot {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slots.Generate(s.availabilites, s.bookings, s.now())
}

// CreateBookingInput is the validated input for CreateBooking. Field-shape
// validation (required fields, email format) happens in the HTTP layer; this
// method enforces the domain rules that need stored state.
type CreateBookingInput struct {
	EventTypeID string
	StartAt     time.Time
	GuestName   string
	GuestEmail  string
}

// CreateBooking enforces the booking invariants and stores the booking
// atomically:
//
//   - the event type must exist            -> ErrEventTypeNotFound
//   - the slot must be bookable by time    -> ErrSlotUnavailable
//   - the slot must not overlap a booking   -> ErrBookingConflict
//
// The existence check, conflict check and insert all run under one lock, so two
// concurrent requests for the same slot cannot both succeed.
func (s *Store) CreateBooking(in CreateBookingInput) (domain.Booking, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.eventTypes[in.EventTypeID]; !ok {
		return domain.Booking{}, ErrEventTypeNotFound
	}

	start := normalize(in.StartAt)
	now := s.now()
	if !slots.IsBookable(s.availabilites, start, now) {
		return domain.Booking{}, ErrSlotUnavailable
	}

	end := start.Add(domain.SlotDuration)
	for _, b := range s.bookings {
		if slots.Overlaps(start, end, b.StartAt, b.EndAt) {
			return domain.Booking{}, ErrBookingConflict
		}
	}

	s.bookingSeq++
	b := domain.Booking{
		ID:          "bkg_" + strconv.Itoa(s.bookingSeq),
		EventTypeID: in.EventTypeID,
		StartAt:     start,
		EndAt:       end,
		GuestName:   strings.TrimSpace(in.GuestName),
		GuestEmail:  strings.TrimSpace(in.GuestEmail),
		CreatedAt:   normalize(now),
	}
	s.bookings = append(s.bookings, b)
	return b, nil
}

// UpcomingBookings returns bookings whose start time is now or later, ordered
// chronologically.
func (s *Store) UpcomingBookings() []domain.Booking {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now().UTC()
	out := make([]domain.Booking, 0, len(s.bookings))
	for _, b := range s.bookings {
		if !b.StartAt.Before(now) {
			out = append(out, b)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].StartAt.Before(out[j].StartAt) })
	return out
}

// normalize converts a time to whole-second UTC so stored values marshal as
// "...:00Z" and compare cleanly.
func normalize(t time.Time) time.Time {
	return t.UTC().Truncate(time.Second)
}
