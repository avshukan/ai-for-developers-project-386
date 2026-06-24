// Command server runs the Call Booking HTTP API.
//
// Configuration (all optional) via environment variables:
//
//	PORT                  TCP port to listen on            (default 8080)
//	CORS_ALLOWED_ORIGIN   Access-Control-Allow-Origin       (default *)
//	SEED_DATA             seed demo data on startup unless  (default true)
//	                      set to "false"
//	STATIC_DIR            directory with the built SPA to   (default empty,
//	                      serve from the same origin         API-only)
//
// Storage is in memory and resets on restart (see
// docs/adr/0002-backend-stack-and-storage.md).
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/avshukan/ai-for-developers-project-386/backend/internal/httpapi"
	"github.com/avshukan/ai-for-developers-project-386/backend/internal/store"
)

func main() {
	addr := ":" + getenv("PORT", "8080")
	allowedOrigin := getenv("CORS_ALLOWED_ORIGIN", "*")
	staticDir := getenv("STATIC_DIR", "")

	st := store.New()
	if getenv("SEED_DATA", "true") != "false" {
		seed(st)
		log.Print("server: seeded demo data (set SEED_DATA=false to disable)")
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           httpapi.NewRouter(st, allowedOrigin, staticDir),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Shut down cleanly on Ctrl-C / SIGTERM.
	go func() {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		<-ctx.Done()
		log.Print("server: shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("server: shutdown error: %v", err)
		}
	}()

	if staticDir != "" {
		log.Printf("server: serving SPA from %s", staticDir)
	}
	log.Printf("server: listening on %s (CORS origin %q)", addr, allowedOrigin)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

// seed loads a couple of event types and a week of availability so the app is
// usable immediately after start. In-memory data resets on restart anyway.
func seed(st *store.Store) {
	st.CreateEventType(store.CreateEventTypeInput{
		Title:           "Intro call",
		Description:     "A 30-minute introduction call.",
		DurationMinutes: 30,
	})
	st.CreateEventType(store.CreateEventTypeInput{
		Title:           "Strategy session",
		Description:     "Discuss your project roadmap.",
		DurationMinutes: 30,
	})

	// Publish 09:00-12:00 and 14:00-17:00 UTC for the next 7 days.
	day := time.Now().UTC().Truncate(24 * time.Hour)
	for i := 1; i <= 7; i++ {
		d := day.AddDate(0, 0, i)
		st.CreateAvailability(store.CreateAvailabilityInput{
			StartAt: d.Add(9 * time.Hour),
			EndAt:   d.Add(12 * time.Hour),
		})
		st.CreateAvailability(store.CreateAvailabilityInput{
			StartAt: d.Add(14 * time.Hour),
			EndAt:   d.Add(17 * time.Hour),
		})
	}
}
