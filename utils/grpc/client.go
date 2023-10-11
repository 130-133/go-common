package grpc

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/zrpc"
	"time"
)

func NewClientWithTimeOut(conf zrpc.RpcClientConf, timeout time.Duration, options ...zrpc.ClientOption) zrpc.Client {
	if timeout.Seconds() == 0 {
		timeout = 10 * time.Millisecond
	}
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	var (
		rpcClient zrpc.Client
		err       error
	)
	for {
		finish := false
		if rpcClient, err = zrpc.NewClient(conf, options...); err == nil {
			break
		}
		select {
		case <-ctx.Done():
			finish = true
		default:
		}
		if finish {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if rpcClient == nil {
		panic(fmt.Sprintf("RPC客户端连接服务器失败, %+v", conf))
	}
	return rpcClient
}
