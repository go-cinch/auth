package server

import (
	"net/http"
)

// DocsHandler serves static files from docs directory.
func DocsHandler() http.Handler {
	// Service runs from cmd/auth, so docs is at ../../docs
	return http.StripPrefix("/docs/", http.FileServer(http.Dir("../../docs")))
}
