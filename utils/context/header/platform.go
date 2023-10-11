package header

import (
	"context"
	"regexp"
	"strings"
)

// Platform 平台
type Platform int

const (
	IOS Platform = iota + 1
	Android
	PC
)

// ExtractPlatform 从UA提取平台
func ExtractPlatform(useragent string) Platform {
	var platform Platform
	useragent = strings.ToLower(useragent)
	ios, _ := regexp.Compile("(iphone|ipad)")
	switch {
	case strings.Contains(useragent, "windows"):
		platform = PC
	case strings.Contains(useragent, "android"):
		platform = Android
	case ios.MatchString(useragent):
		platform = IOS
	}
	return platform
}

// GetPlatformFromCtx 获取平台标识
func GetPlatformFromCtx(ctx context.Context) Platform {
	if u, ok := ctx.Value(CtxUaKey).(UserAgent); ok {
		return u.Platform
	}
	return Platform(0)
}

func (p Platform) IsIOS() bool {
	return p == IOS
}

func (p Platform) IsAndroid() bool {
	return p == Android
}

func (p Platform) IsPC() bool {
	return p == PC
}

func (p Platform) IsOther() bool {
	return !p.IsIOS() && !p.IsAndroid() && !p.IsPC()
}
