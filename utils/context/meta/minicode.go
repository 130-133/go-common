package meta

import (
	"context"
	"reflect"

	"google.golang.org/grpc/metadata"

	"llm-PhotoMagic/go-common/utils/context/header"
)

func ExtractMiniCode(ctx context.Context) header.MiniCode {
	m := header.MiniCode{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return m
	}
	mType := reflect.TypeOf(m)
	mValue := reflect.ValueOf(&m)
	for i := 0; i < mType.NumField(); i++ {
		name := mType.Field(i).Tag.Get("json")
		values := md.Get(name)
		if len(values) == 0 {
			continue
		}
		mValue.Elem().Field(i).SetString(values[0])
	}
	return m
}
