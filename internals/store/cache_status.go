package store

import "context"

type cacheStatusKey struct{}

func ContextWithCacheStatus(parent context.Context) context.Context {
	return context.WithValue(parent, cacheStatusKey{}, new(string))
}

func SetCacheStatus(ctx context.Context, status string) {
	if p, ok := ctx.Value(cacheStatusKey{}).(*string); ok {
		*p = status
	}
}

func GetCacheStatus(ctx context.Context) string {
	if p, ok := ctx.Value(cacheStatusKey{}).(*string); ok {
		return *p
	}
	return ""
}
