package header

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"google.golang.org/grpc/metadata"
)

const CtxMiniCodeKey = "mini-code"

type MiniCode struct {
	Version       string `json:"x-minicode-version"`
	Channel       string `json:"x-minicode-channel"`
	SensorsDataOs string `json:"x-minicode-sensorsdata-os"`
	Platform      string `json:"x-minicode-platform"`
	SystemVersion string `json:"x-minicode-systemversion"`
	Brand         string `json:"x-minicode-brand"`
	Supplier      string `json:"x-minicode-supplier"`
}

func ExtractMiniCode(h http.Header) MiniCode {
	m := MiniCode{}
	mType := reflect.TypeOf(m)
	mValue := reflect.ValueOf(&m)
	for i := 0; i < mType.NumField(); i++ {
		name := mType.Field(i).Tag.Get("json")
		value := h.Get(name)
		mValue.Elem().Field(i).SetString(value)
	}
	return m
}

// GetMiniCodeFromCtx 获取平台标识
func GetMiniCodeFromCtx(ctx context.Context) MiniCode {
	if m, ok := ctx.Value(CtxMiniCodeKey).(MiniCode); ok {
		return m
	}
	return MiniCode{}
}

func (m MiniCode) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, CtxMiniCodeKey, m)
}

func (m MiniCode) InjectMetaData() metadata.MD {
	tmp, _ := json.Marshal(m)
	data := make(map[string]string)
	_ = json.Unmarshal(tmp, &data)
	return metadata.New(data)
}
