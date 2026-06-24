package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/avshukan/ai-for-developers-project-386/backend/internal/domain"
	"github.com/avshukan/ai-for-developers-project-386/backend/internal/httpapi"
	"github.com/avshukan/ai-for-developers-project-386/backend/internal/store"
)

func at(y int, mo time.Month, d, h, mi int) time.Time {
	return time.Date(y, mo, d, h, mi, 0, 0, time.UTC)
}

// newTestServer builds a router over a fixed-clock store with one event type
// and one availability range (2026-06-24 09:00-12:00 UTC -> six slots).
func newTestServer(t *testing.T) (http.Handler, domain.EventType) {
	t.Helper()
	s := store.NewWithClock(func() time.Time { return at(2026, 6, 23, 8, 0) })
	et := s.CreateEventType(store.CreateEventTypeInput{
		Title:           "Intro call",
		Description:     "A 30-minute introduction call.",
		DurationMinutes: 30,
	})
	s.CreateAvailability(store.CreateAvailabilityInput{
		StartAt: at(2026, 6, 24, 9, 0),
		EndAt:   at(2026, 6, 24, 12, 0),
	})
	return httpapi.NewRouter(s, "*", ""), et
}

func do(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body == "" {
		r = httptest.NewRequest(method, path, nil)
	} else {
		r = httptest.NewRequest(method, path, strings.NewReader(body))
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	return rec
}

type errResp struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func decode[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
		t.Fatalf("decode response %q: %v", rec.Body.String(), err)
	}
	return v
}

func requireStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("status = %d, want %d (body: %s)", rec.Code, want, rec.Body.String())
	}
}

func TestListEventTypes(t *testing.T) {
	h, et := newTestServer(t)
	rec := do(t, h, http.MethodGet, "/event-types", "")
	requireStatus(t, rec, http.StatusOK)

	got := decode[[]domain.EventType](t, rec)
	if len(got) != 1 || got[0].ID != et.ID {
		t.Fatalf("event types = %+v, want one with id %q", got, et.ID)
	}
}

func TestCreateEventTypeValidation(t *testing.T) {
	h, _ := newTestServer(t)
	rec := do(t, h, http.MethodPost, "/host/event-types",
		`{"title":"","description":"","durationMinutes":60}`)
	requireStatus(t, rec, http.StatusBadRequest)

	body := decode[errResp](t, rec)
	if body.Code != "validation_error" {
		t.Errorf("code = %q, want validation_error", body.Code)
	}
	if len(body.Details) != 3 {
		t.Errorf("details = %v, want 3 entries", body.Details)
	}
}

func TestCreateEventTypeSuccess(t *testing.T) {
	h, _ := newTestServer(t)
	rec := do(t, h, http.MethodPost, "/host/event-types",
		`{"title":"Strategy","description":"Roadmap chat","durationMinutes":30}`)
	requireStatus(t, rec, http.StatusCreated)

	got := decode[domain.EventType](t, rec)
	if got.ID == "" || got.Title != "Strategy" {
		t.Fatalf("created event type = %+v", got)
	}
}

func TestCreateAvailabilityValidation(t *testing.T) {
	h, _ := newTestServer(t)
	rec := do(t, h, http.MethodPost, "/host/availability",
		`{"startAt":"2026-06-24T12:00:00Z","endAt":"2026-06-24T09:00:00Z"}`)
	requireStatus(t, rec, http.StatusBadRequest)
	if got := decode[errResp](t, rec); got.Code != "validation_error" {
		t.Errorf("code = %q, want validation_error", got.Code)
	}
}

func TestListSlotsNotFound(t *testing.T) {
	h, _ := newTestServer(t)
	rec := do(t, h, http.MethodGet, "/event-types/evt_missing/slots", "")
	requireStatus(t, rec, http.StatusNotFound)
	if got := decode[errResp](t, rec); got.Code != "not_found" {
		t.Errorf("code = %q, want not_found", got.Code)
	}
}

func TestListSlots(t *testing.T) {
	h, et := newTestServer(t)
	rec := do(t, h, http.MethodGet, "/event-types/"+et.ID+"/slots", "")
	requireStatus(t, rec, http.StatusOK)
	if got := decode[[]domain.Slot](t, rec); len(got) != 6 {
		t.Fatalf("got %d slots, want 6", len(got))
	}
}

func TestCreateBookingFlow(t *testing.T) {
	h, et := newTestServer(t)

	rec := do(t, h, http.MethodPost, "/bookings",
		`{"eventTypeId":"`+et.ID+`","startAt":"2026-06-24T09:00:00Z","guestName":"Jane","guestEmail":"jane@example.com"}`)
	requireStatus(t, rec, http.StatusCreated)

	booking := decode[domain.Booking](t, rec)
	if !booking.EndAt.Equal(at(2026, 6, 24, 9, 30)) {
		t.Errorf("endAt = %v, want 09:30", booking.EndAt)
	}

	// The booked slot disappears from the available list.
	slotsRec := do(t, h, http.MethodGet, "/event-types/"+et.ID+"/slots", "")
	if got := decode[[]domain.Slot](t, slotsRec); len(got) != 5 {
		t.Fatalf("got %d slots after booking, want 5", len(got))
	}

	// And it shows up in the host bookings list.
	hostRec := do(t, h, http.MethodGet, "/host/bookings", "")
	requireStatus(t, hostRec, http.StatusOK)
	if got := decode[[]domain.Booking](t, hostRec); len(got) != 1 {
		t.Fatalf("got %d host bookings, want 1", len(got))
	}
}

func TestCreateBookingConflict(t *testing.T) {
	h, et := newTestServer(t)
	body := `{"eventTypeId":"` + et.ID + `","startAt":"2026-06-24T09:00:00Z","guestName":"Jane","guestEmail":"jane@example.com"}`

	if rec := do(t, h, http.MethodPost, "/bookings", body); rec.Code != http.StatusCreated {
		t.Fatalf("first booking status = %d, want 201", rec.Code)
	}

	rec := do(t, h, http.MethodPost, "/bookings", body)
	requireStatus(t, rec, http.StatusConflict)
	if got := decode[errResp](t, rec); got.Code != "booking_conflict" {
		t.Errorf("code = %q, want booking_conflict", got.Code)
	}
}

func TestCreateBookingUnknownEventType(t *testing.T) {
	h, _ := newTestServer(t)
	rec := do(t, h, http.MethodPost, "/bookings",
		`{"eventTypeId":"evt_missing","startAt":"2026-06-24T09:00:00Z","guestName":"Jane","guestEmail":"jane@example.com"}`)
	requireStatus(t, rec, http.StatusNotFound)
	if got := decode[errResp](t, rec); got.Code != "not_found" {
		t.Errorf("code = %q, want not_found", got.Code)
	}
}

func TestCreateBookingInvalidEmail(t *testing.T) {
	h, et := newTestServer(t)
	rec := do(t, h, http.MethodPost, "/bookings",
		`{"eventTypeId":"`+et.ID+`","startAt":"2026-06-24T09:00:00Z","guestName":"Jane","guestEmail":"not-an-email"}`)
	requireStatus(t, rec, http.StatusBadRequest)
	if got := decode[errResp](t, rec); got.Code != "validation_error" {
		t.Errorf("code = %q, want validation_error", got.Code)
	}
}

func TestCreateBookingMalformedJSON(t *testing.T) {
	h, _ := newTestServer(t)
	rec := do(t, h, http.MethodPost, "/bookings", `{not json`)
	requireStatus(t, rec, http.StatusBadRequest)
}

func TestCORSPreflight(t *testing.T) {
	h, _ := newTestServer(t)
	rec := do(t, h, http.MethodOptions, "/bookings", "")
	requireStatus(t, rec, http.StatusNoContent)
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want *", got)
	}
}
