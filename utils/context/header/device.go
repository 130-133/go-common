package header

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"google.golang.org/grpc/metadata"
)

const CtxMiniWorldKey = "mini-world"

type MiniWorld struct {
	Zone       string `json:"x-miniworld-zone"`       //地区
	Channel    string `json:"x-miniworld-channel"`    //渠道
	DeviceId   string `json:"x-miniworld-deviceid"`   //设备ID
	AppVersion string `json:"x-miniworld-appversion"` //APP版本
}

func ExtractMiniWorld(h http.Header) MiniWorld {
	m := MiniWorld{}
	mType := reflect.TypeOf(m)
	mValue := reflect.ValueOf(&m)
	for i := 0; i < mType.NumField(); i++ {
		name := mType.Field(i).Tag.Get("json")
		value := h.Get(name)
		mValue.Elem().Field(i).SetString(value)
	}
	return m
}

// GetMiniWorldFromCtx 获取设备信息
func GetMiniWorldFromCtx(ctx context.Context) MiniWorld {
	if m, ok := ctx.Value(CtxMiniWorldKey).(MiniWorld); ok {
		return m
	}
	return MiniWorld{}
}

func (m MiniWorld) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, CtxMiniWorldKey, m)
}

func (m MiniWorld) InjectMetaData() metadata.MD {
	tmp, _ := json.Marshal(m)
	data := make(map[string]string)
	_ = json.Unmarshal(tmp, &data)
	return metadata.New(data)
}
