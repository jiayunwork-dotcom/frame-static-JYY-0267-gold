package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type errorEnvelope struct {
	OK    bool     `json:"ok"`
	Error APIError `json:"error"`
}

func postSolve(t *testing.T, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/solve", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	SolveHandler(rec, req)
	return rec
}

func TestSolveHandlerOK(t *testing.T) {
	body := `{"nodes":[{"id":"A","x":0,"y":0,"support":"pin"},{"id":"B","x":4,"y":0,"restraint":[false,true,false]}],"elements":[{"from":"A","to":"B","E":2.1e11,"A":0.01,"I":4e-5,"dist":{"q":200}}],"loads":[]}`
	rec := postSolve(t, body)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
	var resp solveResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !resp.OK || len(resp.Reactions) == 0 {
		t.Fatalf("unexpected resp: %+v", resp)
	}
	var totalY float64
	for _, r := range resp.Reactions {
		if r.DOF == "uy" {
			totalY += r.Force
		}
	}
	// Total vertical reaction must balance the applied UDL q*L = 200*4 = 800.
	if totalY < 799 || totalY > 801 {
		t.Fatalf("total vertical reaction = %v, want ~800", totalY)
	}
}

func TestSolveHandlerInvalid(t *testing.T) {
	rec := postSolve(t, `{"nodes":[],"elements":[]}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d", rec.Code)
	}
	var env errorEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if env.Error.Code != "invalid_model" {
		t.Fatalf("code = %q", env.Error.Code)
	}
}

func TestSolveHandlerSingular(t *testing.T) {
	// One member, only ux of A restrained: the body can translate in y and
	// rotate about A, so the reduced stiffness is singular.
	body := `{"nodes":[{"id":"A","x":0,"y":0,"restraint":[true,false,false]},{"id":"B","x":4,"y":0}],"elements":[{"from":"A","to":"B","E":1,"A":1,"I":1}]}`
	rec := postSolve(t, body)
	var env errorEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if env.Error.Code != "singular" {
		t.Fatalf("code = %q body %s", env.Error.Code, rec.Body.String())
	}
}

func TestMetaHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/meta", nil)
	rec := httptest.NewRecorder()
	MetaHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
}
