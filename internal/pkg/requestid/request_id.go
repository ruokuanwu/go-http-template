package requestid

import (
	"context"
	"crypto/rand"
)

func New() string {
	return rand.Text()
}

type contextKey struct{}

func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKey{}, id)
}

func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(contextKey{}).(string)
	return id
}
