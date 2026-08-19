package api

import (
	"errors"
	"fmt"
	"testing"

	"frame-static/internal/assemble"
	"frame-static/internal/model"
)

func TestClassifyInvalidModel(t *testing.T) {
	e := Classify(&model.InvalidModelError{Field: "x", Msg: "bad"})
	if e.Code != "invalid_model" {
		t.Fatalf("code=%q", e.Code)
	}
}

func TestClassifySingular(t *testing.T) {
	e := Classify(assemble.ErrSingular)
	if e.Code != "singular" {
		t.Fatalf("code=%q", e.Code)
	}
}

func TestClassifyWrappedSingular(t *testing.T) {
	wrapped := fmt.Errorf("%w: extra detail", assemble.ErrSingular)
	if !errors.Is(wrapped, assemble.ErrSingular) {
		t.Fatalf("expected wrap detected")
	}
	if got := Classify(wrapped); got.Code != "singular" {
		t.Fatalf("code=%q", got.Code)
	}
}
