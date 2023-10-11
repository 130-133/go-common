package middleware

import (
	"net/http"

	"git.tyy.com/llm-PhotoMagic/go-common/utils/context/header"
)

// UserAgentFun 提取头部信息
func UserAgentFun(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		//提取UA 配合 GetUAFromCtx\GetBrowserFromCtx\GetPlatformFromCtx使用
		ctx = header.ExtractUA(r.Header).WithContext(ctx)
		//提取IP 配合 GetIPFromCtx使用
		ctx = header.ExtractIP(r).WithContext(ctx)
		//提取MINICODE标记 配合 GetMiniCodeFromCtx使用
		ctx = header.ExtractMiniCode(r.Header).WithContext(ctx)
		//提取MINIWORLD标记 配合 GetMiniWorldFromCtx使用
		ctx = header.ExtractMiniWorld(r.Header).WithContext(ctx)
		req := r.WithContext(ctx)
		next(w, req)
	}
}
