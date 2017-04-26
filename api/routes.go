package api

import "github.com/tedsuo/rata"

const (
	HealthCheck = "HealthCheck"
	Webhook     = "Webhook"
)

var Routes = rata.Routes([]rata.Route{
	{Path: "/api/v1/healthcheck", Method: "GET", Name: HealthCheck},
	{Path: "/api/v1/webhook", Method: "POST", Name: Webhook},
})
