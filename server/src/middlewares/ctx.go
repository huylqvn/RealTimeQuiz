package middlewares

import (
	"net/http"
)

type ContextKey string

const (
	ContextKeyReqHeader ContextKey = "Header"
	ContextKeyReqQuery  ContextKey = "Query"
	ContextKeyToken     ContextKey = "Token"
)

func EnrichHeaderCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := EnrichCtx(r.Context(), ContextKeyReqHeader, r.Header)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func EnrichQueryCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := EnrichCtx(r.Context(), ContextKeyReqQuery, r.URL.Query())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
