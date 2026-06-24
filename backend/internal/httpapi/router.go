package httpapi

import (
	"net/http"

	"github.com/avshukan/ai-for-developers-project-386/backend/internal/store"
	"github.com/avshukan/ai-for-developers-project-386/backend/internal/web"
)

// NewRouter wires the API routes from the contract and wraps them with CORS.
//
// allowedOrigin is the value for Access-Control-Allow-Origin (e.g. "*" or the
// SPA's origin). When staticDir is non-empty the built SPA is served from the
// same origin under a catch-all route, so the API and the UI share one port
// (the Docker/Render setup; see docs/adr/0005-deployment-combined-docker-render.md).
// When staticDir is empty the backend is API-only and the SPA is served
// separately (local dev, tests, e2e). The specific API patterns always take
// precedence over the catch-all, so API routing is unaffected.
func NewRouter(s *store.Store, allowedOrigin, staticDir string) http.Handler {
	api := New(s)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /event-types", api.listEventTypes)
	mux.HandleFunc("POST /host/event-types", api.createEventType)
	mux.HandleFunc("POST /host/availability", api.createAvailability)
	mux.HandleFunc("GET /event-types/{eventTypeId}/slots", api.listSlots)
	mux.HandleFunc("POST /bookings", api.createBooking)
	mux.HandleFunc("GET /host/bookings", api.listHostBookings)

	if staticDir != "" {
		mux.Handle("/", web.SPAHandler(staticDir))
	}

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
