package api

import (
	"net/http"

	"github.com/tedsuo/rata"
)

func NewHandler() (http.Handler, error) {
	handlers := map[string]http.Handler{
		HealthCheck: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusOK)
		}),
	}

	return rata.NewRouter(Routes, handlers)
}
