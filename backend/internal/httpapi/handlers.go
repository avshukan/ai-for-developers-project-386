// Package httpapi implements the Call Booking HTTP API defined in
// api/main.tsp. Handlers do field-shape validation and translate store errors
// into the contract's 400 / 404 / 409 responses; the business invariants
// themselves live in the store and slots packages.
package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/avshukan/ai-for-developers-project-386/backend/internal/store"
)

// maxBodyBytes caps request bodies; the MVP payloads are tiny.
const maxBodyBytes = 1 << 20 // 1 MiB

// API holds the dependencies shared by the handlers.
type API struct {
	store *store.Store
}

// New returns an API backed by the given store.
func New(s *store.Store) *API {
	return &API{store: s}
}

type createEventTypeRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	DurationMinutes int    `json:"durationMinutes"`
}

type createAvailabilityRequest struct {
	StartAt time.Time `json:"startAt"`
	EndAt   time.Time `json:"endAt"`
}

type createBookingRequest struct {
	EventTypeID string    `json:"eventTypeId"`
	StartAt     time.Time `json:"startAt"`
	GuestName   string    `json:"guestName"`
	GuestEmail  string    `json:"guestEmail"`
}

// listEventTypes handles GET /event-types.
func (a *API) listEventTypes(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, a.store.ListEventTypes())
}

// createEventType handles POST /host/event-types.
func (a *API) createEventType(w http.ResponseWriter, r *http.Request) {
	var req createEventTypeRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	var details []string
	if isBlank(req.Title) {
		details = append(details, "title must not be empty")
	}
	if isBlank(req.Description) {
		details = append(details, "description must not be empty")
	}
	if req.DurationMinutes != 30 {
		details = append(details, "durationMinutes must be 30")
	}
	if len(details) > 0 {
		writeValidationError(w, "Request validation failed.", details)
		return
	}

	et := a.store.CreateEventType(store.CreateEventTypeInput{
		Title:           req.Title,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
	})
	writeJSON(w, http.StatusCreated, et)
}

// createAvailability handles POST /host/availability.
func (a *API) createAvailability(w http.ResponseWriter, r *http.Request) {
	var req createAvailabilityRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	var details []string
	if req.StartAt.IsZero() {
		details = append(details, "startAt is required")
	}
	if req.EndAt.IsZero() {
		details = append(details, "endAt is required")
	}
	if !req.StartAt.IsZero() && !req.EndAt.IsZero() && !req.EndAt.After(req.StartAt) {
		details = append(details, "endAt must be after startAt")
	}
	if len(details) > 0 {
		writeValidationError(w, "Request validation failed.", details)
		return
	}

	a2 := a.store.CreateAvailability(store.CreateAvailabilityInput{
		StartAt: req.StartAt,
		EndAt:   req.EndAt,
	})
	writeJSON(w, http.StatusCreated, a2)
}

// listSlots handles GET /event-types/{eventTypeId}/slots.
func (a *API) listSlots(w http.ResponseWriter, r *http.Request) {
	eventTypeID := r.PathValue("eventTypeId")
	if !a.store.EventTypeExists(eventTypeID) {
		writeNotFound(w, "Event type not found.")
		return
	}
	writeJSON(w, http.StatusOK, a.store.AvailableSlots())
}

// createBooking handles POST /bookings.
func (a *API) createBooking(w http.ResponseWriter, r *http.Request) {
	var req createBookingRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	var details []string
	if isBlank(req.EventTypeID) {
		details = append(details, "eventTypeId is required")
	}
	if req.StartAt.IsZero() {
		details = append(details, "startAt is required")
	}
	if isBlank(req.GuestName) {
		details = append(details, "guestName is required")
	}
	if isBlank(req.GuestEmail) {
		details = append(details, "guestEmail is required")
	} else if !isValidEmail(req.GuestEmail) {
		details = append(details, "guestEmail must be a valid email")
	}
	if len(details) > 0 {
		writeValidationError(w, "Request validation failed.", details)
		return
	}

	booking, err := a.store.CreateBooking(store.CreateBookingInput{
		EventTypeID: req.EventTypeID,
		StartAt:     req.StartAt,
		GuestName:   req.GuestName,
		GuestEmail:  req.GuestEmail,
	})
	switch {
	case errors.Is(err, store.ErrEventTypeNotFound):
		writeNotFound(w, "Event type not found.")
	case errors.Is(err, store.ErrSlotUnavailable):
		writeValidationError(w, "The selected slot is not available.", []string{
			"selected slot is in the past, outside availability, or not a valid 30-minute slot",
		})
	case errors.Is(err, store.ErrBookingConflict):
		writeConflict(w, "That time was just booked. Please choose another slot.")
	case err != nil:
		writeJSON(w, http.StatusInternalServerError, errorBody{
			Code:    "internal_error",
			Message: "Unexpected error.",
		})
	default:
		writeJSON(w, http.StatusCreated, booking)
	}
}

// listHostBookings handles GET /host/bookings.
func (a *API) listHostBookings(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, a.store.UpcomingBookings())
}

// decodeJSON reads and decodes the JSON request body into dst. On any decode
// failure it writes a 400 validation_error and returns false, so callers can
// simply `if !decodeJSON(...) { return }`.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			writeValidationError(w, "Request body is required.", nil)
		} else {
			writeValidationError(w, "Request body is not valid JSON.", nil)
		}
		return false
	}
	return true
}
