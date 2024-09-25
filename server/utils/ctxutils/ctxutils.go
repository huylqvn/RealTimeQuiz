package ctxutils

import "context"

func EnrichCtx(ctx context.Context, key any, vals interface{}) context.Context {
	return context.WithValue(ctx, key, vals)
}

func GetCtxValue(ctx context.Context, key string) string {
	value := ctx.Value(key)
	if value == nil {
		return ""
	}
	return value.(string)
}
