package meta

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/130-133/go-common/utils/context/header"
)

type Metadata struct {
	metadata.MD
}

func ExtractIp(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	values := md.Get(header.CtxIPKey)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
