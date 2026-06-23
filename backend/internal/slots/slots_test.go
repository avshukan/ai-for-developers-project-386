package slots

import (
	"testing"
	"time"

	"github.com/avshukan/ai-for-developers-project-386/backend/internal/domain"
)

func at(y int, mo time.Month, d, h, mi int) time.Time {
	return time.Date(y, mo, d, h, mi, 0, 0, time.UTC)
}

func avail(start, end time.Time) domain.Availability {
	return domain.Availability{StartAt: start, EndAt: end}
}

func starts(slots []domain.Slot) []time.Time {
	out := make([]time.Time, len(slots))
	for i, s := range slots {
		out[i] = s.StartAt
	}
	return out
}

func equalTimes(a, b []time.Time) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Equal(b[i]) {
			return false
		}
	}
	return true
}

func TestGenerate(t *testing.T) {
	now := at(2026, 6, 23, 8, 0)

	tests := []struct {
		name     string
		avails   []domain.Availability
		bookings []domain.Booking
		now      time.Time
		want     []time.Time
	}{
		{
			name:   "splits range into full 30-minute slots",
			avails: []domain.Availability{avail(at(2026, 6, 24, 9, 0), at(2026, 6, 24, 11, 0))},
			now:    now,
			want: []time.Time{
				at(2026, 6, 24, 9, 0), at(2026, 6, 24, 9, 30),
				at(2026, 6, 24, 10, 0), at(2026, 6, 24, 10, 30),
			},
		},
		{
			name:   "ignores leftover shorter than a slot",
			avails: []domain.Availability{avail(at(2026, 6, 24, 9, 0), at(2026, 6, 24, 10, 20))},
			now:    now,
			want:   []time.Time{at(2026, 6, 24, 9, 0), at(2026, 6, 24, 9, 30)},
		},
		{
			name:   "drops slots that start in the past",
			avails: []domain.Availability{avail(at(2026, 6, 24, 9, 0), at(2026, 6, 24, 11, 0))},
			now:    at(2026, 6, 24, 9, 45),
			want:   []time.Time{at(2026, 6, 24, 10, 0), at(2026, 6, 24, 10, 30)},
		},
		{
			name:   "drops slots outside the 14-day window",
			avails: []domain.Availability{avail(at(2026, 7, 30, 9, 0), at(2026, 7, 30, 11, 0))},
			now:    now,
			want:   nil,
		},
		{
			name:     "excludes slots overlapping a booking",
			avails:   []domain.Availability{avail(at(2026, 6, 24, 9, 0), at(2026, 6, 24, 11, 0))},
			bookings: []domain.Booking{{StartAt: at(2026, 6, 24, 9, 30), EndAt: at(2026, 6, 24, 10, 0)}},
			now:      now,
			want: []time.Time{
				at(2026, 6, 24, 9, 0), at(2026, 6, 24, 10, 0), at(2026, 6, 24, 10, 30),
			},
		},
		{
			name: "de-duplicates and sorts slots from overlapping ranges",
			avails: []domain.Availability{
				avail(at(2026, 6, 24, 9, 30), at(2026, 6, 24, 10, 30)),
				avail(at(2026, 6, 24, 9, 0), at(2026, 6, 24, 10, 0)),
			},
			now:  now,
			want: []time.Time{at(2026, 6, 24, 9, 0), at(2026, 6, 24, 9, 30), at(2026, 6, 24, 10, 0)},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Generate(tc.avails, tc.bookings, tc.now)
			if !equalTimes(starts(got), tc.want) {
				t.Fatalf("slot starts = %v, want %v", starts(got), tc.want)
			}
			for _, s := range got {
				if s.Status != domain.SlotAvailable {
					t.Errorf("slot %v status = %q, want available", s.StartAt, s.Status)
				}
				if !s.EndAt.Equal(s.StartAt.Add(domain.SlotDuration)) {
					t.Errorf("slot %v endAt = %v, want start+30m", s.StartAt, s.EndAt)
				}
			}
		})
	}
}

func TestIsBookable(t *testing.T) {
	now := at(2026, 6, 23, 8, 0)
	avails := []domain.Availability{avail(at(2026, 6, 24, 9, 0), at(2026, 6, 24, 11, 0))}

	tests := []struct {
		name  string
		start time.Time
		want  bool
	}{
		{"aligned slot inside range", at(2026, 6, 24, 9, 0), true},
		{"last full slot inside range", at(2026, 6, 24, 10, 30), true},
		{"misaligned start", at(2026, 6, 24, 9, 15), false},
		{"start before range", at(2026, 6, 24, 8, 30), false},
		{"slot extends past range end", at(2026, 6, 24, 11, 0), false},
		{"in the past", at(2026, 6, 20, 9, 0), false},
		{"outside 14-day window", at(2026, 7, 30, 9, 0), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsBookable(avails, tc.start, now); got != tc.want {
				t.Fatalf("IsBookable(%v) = %v, want %v", tc.start, got, tc.want)
			}
		})
	}
}

func TestOverlaps(t *testing.T) {
	a, b := at(2026, 6, 24, 9, 0), at(2026, 6, 24, 9, 30)
	c, d := at(2026, 6, 24, 9, 30), at(2026, 6, 24, 10, 0)

	if Overlaps(a, b, c, d) {
		t.Error("adjacent half-open intervals must not overlap")
	}
	if !Overlaps(a, b, at(2026, 6, 24, 9, 15), at(2026, 6, 24, 9, 45)) {
		t.Error("intersecting intervals must overlap")
	}
}
