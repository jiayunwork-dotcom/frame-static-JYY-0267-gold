package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// InvalidModelError describes a structural defect that prevents solving.
type InvalidModelError struct {
	Field string
	Msg   string
}

func (e *InvalidModelError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Msg)
	}
	return e.Msg
}

// ParseModel decodes a JSON model from r.
func ParseModel(r io.Reader) (*Model, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var m Model
	if err := dec.Decode(&m); err != nil {
		return nil, &InvalidModelError{Field: "json", Msg: err.Error()}
	}
	return &m, nil
}

// ParseModelBytes is a convenience wrapper around ParseModel.
func ParseModelBytes(b []byte) (*Model, error) {
	return ParseModel(bytes.NewReader(b))
}
