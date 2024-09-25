package middlewares

import (
	"context"
	"net/http"
	"quizserver/src/endpoints/auth"
	"quizserver/src/errs"
	"quizserver/src/service"
	"strings"

	"github.com/go-kit/kit/endpoint"
)

func EnrichCtx(ctx context.Context, key any, vals interface{}) context.Context {
	return context.WithValue(ctx, key, vals)
}

func Auth(s *service.Service, f func(s *service.Service) endpoint.Endpoint) endpoint.Endpoint {
	return endpoint.Chain(authentication(s))(f(s))
}

func authentication(s *service.Service) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			token := tokenLookup(ctx)
			if token == "" {
				return nil, errs.WrapMessage(errs.ErrAuthorization, "token not found")
			}

			claims, err := auth.ValidateJwtToken(s.Config.JwtSecret, token)
			if err != nil {
				return nil, errs.WrapMessage(errs.ErrAuthorization, err.Error())
			}

			ctx = EnrichCtx(ctx, "UserID", claims.UserID)
			return next(ctx, request)
		}
	}
}

func tokenLookup(ctx context.Context) string {
	header, yes := ctx.Value(ContextKeyReqHeader).(http.Header)
	if !yes {
		return ""
	}

	authorization := header.Get("Authorization")
	bearerToken := strings.Split(authorization, " ")
	if len(bearerToken) > 1 {
		return bearerToken[1]
	}
	return ""
}
