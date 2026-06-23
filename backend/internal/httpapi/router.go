package httpapi

import (
	"net/http"

	"github.com/avshukan/ai-for-developers-project-386/backend/internal/store"
)

// NewRouter wires the API routes from the contract and wraps them with CORS so
// the separately-served SPA can call the backend cross-origin. allowedOrigin is
// the value for Access-Control-Allow-Origin (e.g. "*" or the SPA's origin).
func NewRouter(s *store.Store, allowedOrigin string) http.Handler {
	api := New(s)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /event-types", api.listEventTypes)
	mux.HandleFunc("POST /host/event-types", api.createEventType)
	mux.HandleFunc("POST /host/availability", api.createAvailability)
	mux.HandleFunc("GET /event-types/{eventTypeId}/slots", api.listSlots)
	mux.HandleFunc("POST /bookings", api.createBooking)
	mux.HandleFunc("GET /host/bookings", api.listHostBookings)

	return withCORS(mux, allowedOrigin)
}

// withCORS adds permissive CORS headers and short-circuits preflight requests.
// Host pages are public in the MVP (docs/domain.md), so a single allowed origin
// is enough; there are no credentials to protect.
func withCORS(next http.Handler, allowedOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
