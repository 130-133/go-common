package auth

import (
	"context"
)

const CtxKey = "AuthNormal"

func ExtractJwt(ctx context.Context) map[string]interface{} {
	data := ctx.Value(CtxKey).(map[string]interface{})
	if data == nil {
		return make(map[string]interface{})
	}
	return data
}

func InjectJwt(ctx context.Context, data map[string]interface{}) context.Context {
	return context.WithValue(ctx, CtxKey, data)
}
