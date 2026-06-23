package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
)

// Error codes from the API contract (api/main.tsp).
const (
	codeValidation = "validation_error"
	codeNotFound   = "not_found"
	codeConflict   = "booking_conflict"
)

// errorBody is the shape of the 400 / 404 / 409 error responses in the contract.
type errorBody struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// writeJSON serializes v as JSON with the given status code. Encoding happens
// after the status is set; an encode failure can only be logged at that point.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("httpapi: failed to encode response: %v", err)
	}
}

// writeValidationError returns 400 validation_error with optional field details.
func writeValidationError(w http.ResponseWriter, message string, details []string) {
	writeJSON(w, http.StatusBadRequest, errorBody{
		Code:    codeValidation,
		Message: message,
		Details: details,
	})
}

// writeNotFound returns 404 not_found.
func writeNotFound(w http.ResponseWriter, message string) {
	writeJSON(w, http.StatusNotFound, errorBody{Code: codeNotFound, Message: message})
}

// writeConflict returns 409 booking_conflict.
func writeConflict(w http.ResponseWriter, message string) {
	writeJSON(w, http.StatusConflict, errorBody{Code: codeConflict, Message: message})
}
