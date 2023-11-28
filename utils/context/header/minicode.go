package header

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"google.golang.org/grpc/metadata"
)

const CtxCodeKey = "code"

type Code struct {
	Version       string `json:"x-version"`
	Channel       string `json:"x-channel"`
	Platform      string `json:"x-platform"`
	SystemVersion string `json:"x-systemVersion"`
	Brand         string `json:"x-brand"`
	Supplier      string `json:"x-supplier"`
}

func ExtractCode(h http.Header) Code {
	m := Code{}
	mType := reflect.TypeOf(m)
	mValue := reflect.ValueOf(&m)
	for i := 0; i < mType.NumField(); i++ {
		name := mType.Field(i).Tag.Get("json")
		value := h.Get(name)
		mValue.Elem().Field(i).SetString(value)
	}
	return m
}

// GetCodeFromCtx 获取平台标识
func GetCodeFromCtx(ctx context.Context) Code {
	if m, ok := ctx.Value(CtxCodeKey).(Code); ok {
		return m
	}
	return Code{}
}

func (m Code) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, CtxCodeKey, m)
}

func (m Code) InjectMetaData() metadata.MD {
	tmp, _ := json.Marshal(m)
	data := make(map[string]string)
	_ = json.Unmarshal(tmp, &data)
	return metadata.New(data)
}
