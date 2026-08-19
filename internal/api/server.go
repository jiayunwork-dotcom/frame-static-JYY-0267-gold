package api

import (
	"io/fs"
	"net/http"
)

// New builds the HTTP handler: API routes plus static web/example serving.
func New(webFS, exampleFS fs.FS) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/solve", SolveHandler)
	mux.HandleFunc("/api/meta", MetaHandler)
	mux.Handle("/", staticHandler(webFS, exampleFS))
	return withLogging(mux)
}
