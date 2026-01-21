package service

import (
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel"

	"github.com/go-cinch/common/log"
)

type HealthStatus struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks,omitempty"`
}

// HealthCheck handles the HTTP health endpoint.
func (s *AuthService) HealthCheck(writer http.ResponseWriter, req *http.Request) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(req.Context(), "HealthCheck")
	defer span.End()

	status := HealthStatus{
		Status: "healthy",
		Checks: make(map[string]string),
	}
	httpStatus := http.StatusOK

	// Check DB
	if err := s.health.PingDB(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("health check: db ping failed")
		status.Status = "unhealthy"
		status.Checks["db"] = err.Error()
		httpStatus = http.StatusServiceUnavailable
	} else {
		status.Checks["db"] = "ok"
	}

	// Check Redis
	if err := s.health.PingRedis(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("health check: redis ping failed")
		status.Status = "unhealthy"
		status.Checks["redis"] = err.Error()
		httpStatus = http.StatusServiceUnavailable
	} else {
		status.Checks["redis"] = "ok"
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpStatus)
	json.NewEncoder(writer).Encode(status)
}
