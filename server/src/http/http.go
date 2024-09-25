package http

import (
	"context"
	"net/http"
	"quizserver/src/endpoints"
	"quizserver/src/middlewares"
	"quizserver/src/service"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.elastic.co/apm/module/apmhttp/v2"
)

//	@title			Quiz API
//	@version		1.0.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/api/v1

//	@securityDefinitions.apikey	JWTHeader
//	@in							header
//	@name						Authorization
//	@description				JWT Bearer for access with format "Bearer [token]"

//	@securityDefinitions.apikey	JWTQuery
//	@in							query
//	@name						access_token
//	@description				JWT Bearer for access in query param for get file url request

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

func NewHTTPHandler(s service.Service, endpoints endpoints.Endpoints, logger log.Logger, useCORS bool) http.Handler {
	r := chi.NewRouter()

	origins := []string{"*"}
	if useCORS {
		cors := cors.New(cors.Options{
			AllowedOrigins: origins,
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
				http.MethodOptions,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: false,
			MaxAge:           600,
		})
		r.Use(cors.Handler)
	}

	r.Use(middleware.Recoverer)

	defaultCompressibleContentTypes := []string{
		"text/html",
		"text/css",
		"text/plain",
		"text/javascript",
		"application/javascript",
		"application/x-javascript",
		"application/json",
		"application/atom+xml",
		"application/rss+xml",
		"image/svg+xml",
	}
	r.Use(middleware.Compress(5, defaultCompressibleContentTypes...))

	r.Use(httprate.LimitByIP(100, 30*time.Second))

	r.Use(middlewares.EnrichHeaderCtx)
	r.Use(middlewares.EnrichQueryCtx)

	if s.Config.Env == "development" {
		r.Get("/api/swagger/*", httpSwagger.Handler())
	}

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}
	initRoute(r, endpoints, options)
	initAuthRoute(r, &s)

	return apmhttp.Wrap(r)
}

func DecodeNullRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func initAuthRoute(r chi.Router, s *service.Service) {

}

func initRoute(r chi.Router, endpoints endpoints.Endpoints, options []httptransport.ServerOption) {
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", httptransport.NewServer(
			endpoints.HealthCheck,
			DecodeNullRequest,
			encodeJSONResponse,
			options...,
		).ServeHTTP)
	})
}
