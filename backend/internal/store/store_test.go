package store

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/avshukan/ai-for-developers-project-386/backend/internal/domain"
)

func at(y int, mo time.Month, d, h, mi int) time.Time {
	return time.Date(y, mo, d, h, mi, 0, 0, time.UTC)
}

// fixtureStore returns a store pinned to a fixed clock with one event type and
// one availability range (2026-06-24 09:00-12:00 UTC -> six 30-minute slots).
func fixtureStore(t *testing.T) (*Store, domain.EventType) {
	t.Helper()
	s := NewWithClock(func() time.Time { return at(2026, 6, 23, 8, 0) })
	et := s.CreateEventType(CreateEventTypeInput{
		Title:           "Intro call",
		Description:     "A 30-minute introduction call.",
		DurationMinutes: 30,
	})
	s.CreateAvailability(CreateAvailabilityInput{
		StartAt: at(2026, 6, 24, 9, 0),
		EndAt:   at(2026, 6, 24, 12, 0),
	})
	return s, et
}

func TestCreateEventTypeAndList(t *testing.T) {
	s := New()
	et := s.CreateEventType(CreateEventTypeInput{Title: "Intro", Description: "Desc", DurationMinutes: 30})

	if et.ID == "" {
		t.Fatal("expected a generated id")
	}
	if !s.EventTypeExists(et.ID) {
		t.Errorf("EventTypeExists(%q) = false, want true", et.ID)
	}
	if got := s.ListEventTypes(); len(got) != 1 || got[0].ID != et.ID {
		t.Errorf("ListEventTypes() = %+v, want one entry with id %q", got, et.ID)
	}
}

func TestCreateBookingHappyPath(t *testing.T) {
	s, et := fixtureStore(t)

	b, err := s.CreateBooking(CreateBookingInput{
		EventTypeID: et.ID,
		StartAt:     at(2026, 6, 24, 9, 0),
		GuestName:   "Jane Doe",
		GuestEmail:  "jane@example.com",
	})
	if err != nil {
		t.Fatalf("CreateBooking returned error: %v", err)
	}
	if !b.EndAt.Equal(at(2026, 6, 24, 9, 30)) {
		t.Errorf("endAt = %v, want start + 30m", b.EndAt)
	}

	// The booked slot must no longer be offered as available.
	for _, slot := range s.AvailableSlots() {
		if slot.StartAt.Equal(at(2026, 6, 24, 9, 0)) {
			t.Error("booked 09:00 slot is still listed as available")
		}
	}
}

func TestCreateBookingUnknownEventType(t *testing.T) {
	s, _ := fixtureStore(t)

	_, err := s.CreateBooking(CreateBookingInput{
		EventTypeID: "evt_missing",
		StartAt:     at(2026, 6, 24, 9, 0),
		GuestName:   "Jane",
		GuestEmail:  "jane@example.com",
	})
	if !errors.Is(err, ErrEventTypeNotFound) {
		t.Fatalf("err = %v, want ErrEventTypeNotFound", err)
	}
}

func TestCreateBookingSlotUnavailable(t *testing.T) {
	s, et := fixtureStore(t)

	cases := map[string]time.Time{
		"in the past":          at(2026, 6, 20, 9, 0),
		"outside availability": at(2026, 6, 24, 13, 0),
		"misaligned":           at(2026, 6, 24, 9, 15),
	}
	for name, start := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := s.CreateBooking(CreateBookingInput{
				EventTypeID: et.ID, StartAt: start, GuestName: "Jane", GuestEmail: "jane@example.com",
			})
			if !errors.Is(err, ErrSlotUnavailable) {
				t.Fatalf("err = %v, want ErrSlotUnavailable", err)
			}
		})
	}
}

func TestCreateBookingConflictAcrossEventTypes(t *testing.T) {
	s, et1 := fixtureStore(t)
	et2 := s.CreateEventType(CreateEventTypeInput{Title: "Strategy", Description: "Desc", DurationMinutes: 30})

	if _, err := s.CreateBooking(CreateBookingInput{
		EventTypeID: et1.ID, StartAt: at(2026, 6, 24, 9, 0), GuestName: "Jane", GuestEmail: "jane@example.com",
	}); err != nil {
		t.Fatalf("first booking failed: %v", err)
	}

	// Same time, different event type: must still be rejected (global rule).
	_, err := s.CreateBooking(CreateBookingInput{
		EventTypeID: et2.ID, StartAt: at(2026, 6, 24, 9, 0), GuestName: "Sam", GuestEmail: "sam@example.com",
	})
	if !errors.Is(err, ErrBookingConflict) {
		t.Fatalf("err = %v, want ErrBookingConflict", err)
	}
}

func TestUpcomingBookingsFiltersAndSorts(t *testing.T) {
	// Clock set so that one booked slot is already in the past.
	s := NewWithClock(func() time.Time { return at(2026, 6, 24, 10, 15) })
	et := s.CreateEventType(CreateEventTypeInput{Title: "Intro", Description: "Desc", DurationMinutes: 30})
	s.CreateAvailability(CreateAvailabilityInput{StartAt: at(2026, 6, 24, 9, 0), EndAt: at(2026, 6, 24, 12, 0)})

	// Insert directly to bypass the "no past booking" rule for setup.
	s.bookings = []domain.Booking{
		{ID: "bkg_past", StartAt: at(2026, 6, 24, 9, 0), EndAt: at(2026, 6, 24, 9, 30)},
		{ID: "bkg_late", StartAt: at(2026, 6, 24, 11, 0), EndAt: at(2026, 6, 24, 11, 30)},
		{ID: "bkg_soon", StartAt: at(2026, 6, 24, 10, 30), EndAt: at(2026, 6, 24, 11, 0)},
	}
	_ = et

	got := s.UpcomingBookings()
	if len(got) != 2 {
		t.Fatalf("got %d upcoming bookings, want 2", len(got))
	}
	if got[0].ID != "bkg_soon" || got[1].ID != "bkg_late" {
		t.Errorf("order = [%s %s], want [bkg_soon bkg_late]", got[0].ID, got[1].ID)
	}
}

// TestCreateBookingConcurrent verifies the no-double-booking invariant under
// concurrent requests for the same slot: exactly one must succeed. Run with
// -race to catch data races on the shared store.
func TestCreateBookingConcurrent(t *testing.T) {
	s, et := fixtureStore(t)

	const n = 50
	var wg sync.WaitGroup
	results := make([]error, n)
	start := make(chan struct{})

	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_, results[i] = s.CreateBooking(CreateBookingInput{
				EventTypeID: et.ID,
				StartAt:     at(2026, 6, 24, 9, 0),
				GuestName:   "Guest",
				GuestEmail:  "guest@example.com",
			})
		}(i)
	}
	close(start)
	wg.Wait()

	var ok, conflicts int
	for _, err := range results {
		switch {
		case err == nil:
			ok++
		case errors.Is(err, ErrBookingConflict):
			conflicts++
		default:
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if ok != 1 || conflicts != n-1 {
		t.Fatalf("got %d ok / %d conflicts, want 1 / %d", ok, conflicts, n-1)
	}
}
