package auth

import "context"

func GetId(ctx context.Context) int64 {
	val, _ := ExtractJwt(ctx)["id"].(int64)
	return val
}

func GetUin(ctx context.Context) string {
	val, _ := ExtractJwt(ctx)["uin"].(string)
	return val
}

func GetUserId(ctx context.Context) int64 {
	val, _ := ExtractJwt(ctx)["user_id"].(float64)
	return int64(val)
}

func GetSignedUnix(ctx context.Context) int64 {
	val, _ := ExtractJwt(ctx)["signed_unix"].(float64)
	return int64(val)
}
