// Package web serves the built Vue SPA from the same origin as the API.
//
// It is only mounted when a static directory is configured (see
// httpapi.NewRouter and the STATIC_DIR env var). Local development, unit tests,
// and the Playwright e2e suite run the backend API-only and let Vite serve the
// SPA separately; the combined single-origin setup is used in the Docker image
// that ships to Render (see docs/adr/0005-deployment-combined-docker-render.md).
package web

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// SPAHandler serves a single-page application from dir: it returns the requested
// file when it exists and falls back to index.html for every other path so the
// client-side router can resolve deep links (e.g. a refresh on /bookings).
func SPAHandler(dir string) http.Handler {
	root := filepath.Clean(dir)
	indexPath := filepath.Join(root, "index.html")
	fileServer := http.FileServer(http.Dir(root))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Map the URL path onto a file inside root, guarding against traversal.
		target := filepath.Join(root, filepath.FromSlash(filepath.Clean("/"+r.URL.Path)))
		if target != root && !strings.HasPrefix(target, root+string(os.PathSeparator)) {
			http.NotFound(w, r)
			return
		}

		if info, err := os.Stat(target); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Unknown path or a directory: serve the SPA entry point.
		http.ServeFile(w, r, indexPath)
	})
}
