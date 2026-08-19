package api

import (
	"io"
	"net/http"

	"frame-static/internal/assemble"
	"frame-static/internal/model"
)

const maxBody = 1 << 20 // 1 MiB

// SolveHandler accepts a JSON model and returns the solved displacements,
// member end forces, and support reactions.
func SolveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, APIError{Code: "method", Message: "POST required"})
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBody))
	if err != nil {
		writeError(w, APIError{Code: "read", Message: err.Error()})
		return
	}
	m, err := model.ParseModelBytes(body)
	if err != nil {
		writeError(w, Classify(err))
		return
	}
	res, err := assemble.Solve(m)
	if err != nil {
		writeError(w, Classify(err))
		return
	}
	writeJSON(w, buildResponse(res))
}

// MetaHandler reports the service capabilities.
func MetaHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]any{
		"service":  "frame-static",
		"method":   "direct stiffness (reduced system)",
		"endpoint": "/api/solve",
	})
}
