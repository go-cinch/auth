package server

import (
	"net/http"

	"auth/internal/service"
)

func HealthHandler(svc *service.AuthService) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/pub/healthcheck", svc.HealthCheck)
	return mux
}
