package header

import (
	"context"
	"net/http"
)

const CtxUaKey = "ua"

type UserAgent struct {
	UA       string   `json:"ua"`
	Browser  Browser  `json:"browser"`
	Platform Platform `json:"platform"`
	Language string   `json:"language"`
}

func ExtractUA(h http.Header) UserAgent {
	useragent := h.Get("user-agent")
	return UserAgent{
		UA:       useragent,
		Browser:  ExtractBrowser(useragent),
		Platform: ExtractPlatform(useragent),
		Language: h.Get("User-Language"),
	}
}

// GetUAFromCtx 获取UA
func GetUAFromCtx(ctx context.Context) string {
	if u, ok := ctx.Value(CtxUaKey).(UserAgent); ok {
		return u.UA
	}
	return ""
}

func (u UserAgent) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, CtxUaKey, u)
}
