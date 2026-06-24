# syntax=docker/dockerfile:1
#
# Single combined image for the Call Booking app: the Go API also serves the
# built Vue SPA from the same origin on one port ($PORT), so the whole app is one
# container with one public URL. See docs/adr/0005-deployment-combined-docker-render.md.

# --- Stage 1: build the SPA -------------------------------------------------
FROM node:22-alpine AS frontend
WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ ./
# Empty base URL => the SPA calls the API with same-origin relative paths
# (see frontend/src/api/client.ts). Overridable for a split deployment.
ARG VITE_API_BASE_URL=""
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL
RUN npm run build   # -> /app/frontend/dist

# --- Stage 2: build the Go server ------------------------------------------
FROM golang:1.26-alpine AS backend
WORKDIR /app/backend

COPY backend/go.mod ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# --- Stage 3: minimal runtime ----------------------------------------------
FROM alpine:3.20
WORKDIR /app

# Run as an unprivileged user.
RUN adduser -D -H app
USER app

COPY --from=backend /server /app/server
COPY --from=frontend /app/frontend/dist /app/web

# PORT is injected by the platform (default 8080 in the server); STATIC_DIR turns
# on same-origin SPA serving; SEED_DATA loads demo data on every start.
ENV STATIC_DIR=/app/web \
    SEED_DATA=true
EXPOSE 8080

CMD ["/app/server"]
