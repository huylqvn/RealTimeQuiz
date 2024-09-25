package endpoints

import (
	"github.com/go-kit/kit/endpoint"

	"quizserver/src/endpoints/app"
	"quizserver/src/service"
)

type Endpoints struct {
	HealthCheck endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct
func MakeServerEndpoints(s *service.Service) Endpoints {
	return Endpoints{
		HealthCheck: app.HealthCheckHandler(s),
	}
}
