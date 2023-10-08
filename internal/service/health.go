package service

import (
	"net/http"

	"github.com/go-cinch/common/log"
)

func (*AuthService) HealthCheck(writer http.ResponseWriter, _ *http.Request) {
	log.Info("healthcheck")
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("{}"))
	return
}
