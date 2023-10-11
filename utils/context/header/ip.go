package header

import (
	"context"
	"net"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

const CtxIPKey = "mini-ip"

type ClientIP struct {
	net.IP
}

func ExtractIP(r *http.Request) ClientIP {
	forwarded := r.Header.Get("x-forwarded-for")
	realIP := r.Header.Get("x-real-ip")

	forwardedArr := strings.Split(forwarded, ",")
	if len(forwardedArr) > 0 {
		realIP = forwardedArr[0]
	}
	if realIP == "" {
		realIP, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	ip := net.ParseIP(realIP)
	return ClientIP{
		ip,
	}
}

// GetIPFromCtx 获取浏览器标识
func GetIPFromCtx(ctx context.Context) string {
	if i, ok := ctx.Value(CtxIPKey).(ClientIP); ok {
		if i.String() == "<nil>" {
			return ""
		}
		return i.String()
	}
	return ""
}

func (ip ClientIP) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, CtxIPKey, ip)
}

func (ip ClientIP) InjectMetaData() metadata.MD {
	return metadata.Pairs(CtxIPKey, ip.String())
}
