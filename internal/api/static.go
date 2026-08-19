package api

import (
	"io/fs"
	"net/http"
)

// staticHandler serves the embedded web directory at "/" and the embedded
// example files at "/example/".
func staticHandler(webFS, exampleFS fs.FS) http.Handler {
	mux := http.NewServeMux()
	idx, err := fs.Sub(webFS, ".")
	if err != nil {
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.FS(idx)))
	ex, err := fs.Sub(exampleFS, ".")
	if err != nil {
		panic(err)
	}
	mux.Handle("/example/", http.StripPrefix("/example/", http.FileServer(http.FS(ex))))
	return mux
}
