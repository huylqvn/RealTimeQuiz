package app

import (
	"context"
	"quizserver/src/domain"
	"quizserver/src/service"
	"time"

	"github.com/go-kit/kit/endpoint"
	"go.elastic.co/apm/v2"
)

func HealthCheckHandler(s *service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		startTime := time.Now()
		span, ctx := apm.StartSpan(ctx, "HealthCheckHandler", "endpoint")
		defer func() {
			s.Logger.Trace("HealthCheckHandler", startTime, nil, nil, span)
		}()
		return domain.SuccessResponse{Message: "server is running"}, nil
	}
}
