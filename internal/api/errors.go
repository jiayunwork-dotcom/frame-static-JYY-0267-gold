package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"frame-static/internal/assemble"
	"frame-static/internal/model"
)

// APIError is a machine-readable error returned to clients.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string { return e.Message }

// Classify maps a solver error to an APIError with a stable code so clients can
// branch on the cause without parsing the message.
func Classify(err error) APIError {
	var inv *model.InvalidModelError
	if errors.As(err, &inv) {
		return APIError{Code: "invalid_model", Message: err.Error()}
	}
	if errors.Is(err, assemble.ErrSingular) {
		return APIError{Code: "singular", Message: "structure is a mechanism or under-constrained: no unique solution"}
	}
	return APIError{Code: "solve_failed", Message: err.Error()}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "encode failed", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, e APIError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": false, "error": e})
}
