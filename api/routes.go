package api

import "github.com/tedsuo/rata"

const (
	HealthCheck = "HealthCheck"
)

var Routes = rata.Routes([]rata.Route{
	{Path: "/api/v1/healthcheck", Method: "GET", Name: HealthCheck},
})