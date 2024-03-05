package middleware_grpc

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"

	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/logger"
)

type IGrpcLog interface {
	Interceptor() grpc.UnaryServerInterceptor
}

type log struct {
	*logger.MLogger
	ignore []string
}

// GrpcLog rpc请求日志
func GrpcLog(l *logger.MLogger) IGrpcLog {
	return log{
		MLogger: l,
		ignore: []string{
			"/ping",
			"/checkhealth",
		},
	}
}

func (l log) Interceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if l.ignorePath(info.FullMethod) {
			return handler(ctx, req)
		}
		l.WithCtx(ctx).WithReq(req).Info(fmt.Sprintf("请求参数 - %s", info.FullMethod))
		return handler(ctx, req)
	}
}

func (l log) ignorePath(path string) bool {
	for _, v := range l.ignore {
		if strings.HasSuffix(strings.ToLower(path), v) {
			return true
		}
	}
	return false
}
