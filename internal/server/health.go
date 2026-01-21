package server

import (
	"net/http"

	"auth/internal/service"
)

// HealthHandler exposes the health endpoints served by the HTTP mux.
func HealthHandler(svc *service.AuthService) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", svc.HealthCheck)
	return mux
}
