// Package slots derives bookable 30-minute slots from host availability.
//
// All functions are pure: they take availability ranges, existing bookings and
// the current time, and return results without touching shared state. This
// keeps the slot rules from docs/domain.md ("Slot Generation Rules") in one
// place and easy to test in isolation.
package slots

import (
	"sort"
	"time"

	"github.com/avshukan/ai-for-developers-project-386/backend/internal/domain"
)

// Generate returns the available slots within the next 14 days, in
// chronological order. A 30-minute interval becomes an available slot when it
// fits fully inside an availability range, does not start in the past, starts
// within the slot window, and does not overlap any existing booking.
// Slots produced by overlapping availability ranges are de-duplicated by start
// time.
func Generate(avails []domain.Availability, bookings []domain.Booking, now time.Time) []domain.Slot {
	now = now.UTC()
	windowEnd := now.Add(domain.SlotWindow)

	seen := make(map[int64]struct{})
	out := make([]domain.Slot, 0)

	for _, a := range avails {
		for _, iv := range splitRange(a.StartAt, a.EndAt) {
			if iv.start.Before(now) || !iv.start.Before(windowEnd) {
				continue // past, or outside the 14-day window
			}
			if overlapsAnyBooking(iv.start, iv.end, bookings) {
				continue // booked time is not available
			}
			key := iv.start.Unix()
			if _, dup := seen[key]; dup {
				continue
			}
			seen[key] = struct{}{}
			out = append(out, domain.Slot{
				StartAt: iv.start,
				EndAt:   iv.end,
				Status:  domain.SlotAvailable,
			})
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i].StartAt.Before(out[j].StartAt) })
	return out
}

// IsBookable reports whether startAt is the start of a valid, bookable slot:
// aligned to an availability range's 30-minute grid, fully inside that range,
// not in the past, and within the slot window. It deliberately ignores
// existing bookings — conflict detection is a separate, atomic step so that a
// real-but-taken slot can be reported as a 409 conflict rather than a 400.
func IsBookable(avails []domain.Availability, startAt, now time.Time) bool {
	startAt = startAt.UTC()
	now = now.UTC()
	end := startAt.Add(domain.SlotDuration)

	if startAt.Before(now) || !startAt.Before(now.Add(domain.SlotWindow)) {
		return false
	}

	for _, a := range avails {
		offset := startAt.Sub(a.StartAt.UTC())
		if offset < 0 || offset%domain.SlotDuration != 0 {
			continue // not aligned to this range's grid
		}
		if !end.After(a.EndAt.UTC()) {
			return true // fits fully inside the range
		}
	}
	return false
}

// Overlaps reports whether two half-open time intervals [aStart, aEnd) and
// [bStart, bEnd) intersect. This is the no-double-booking check from
// docs/domain.md: new.startAt < existing.endAt AND new.endAt > existing.startAt.
func Overlaps(aStart, aEnd, bStart, bEnd time.Time) bool {
	return aStart.Before(bEnd) && aEnd.After(bStart)
}

type interval struct {
	start time.Time
	end   time.Time
}

// splitRange splits [start, end) into consecutive full 30-minute intervals,
// ignoring any leftover shorter than a full slot.
func splitRange(start, end time.Time) []interval {
	start = start.UTC()
	end = end.UTC()
	var ivs []interval
	for cursor := start; !cursor.Add(domain.SlotDuration).After(end); cursor = cursor.Add(domain.SlotDuration) {
		ivs = append(ivs, interval{start: cursor, end: cursor.Add(domain.SlotDuration)})
	}
	return ivs
}

func overlapsAnyBooking(start, end time.Time, bookings []domain.Booking) bool {
	for _, b := range bookings {
		if Overlaps(start, end, b.StartAt.UTC(), b.EndAt.UTC()) {
			return true
		}
	}
	return false
}
