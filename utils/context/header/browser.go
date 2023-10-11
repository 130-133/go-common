package header

import (
	"context"
	"strings"
)

// Browser 浏览器
type Browser int

const (
	Wechat Browser = iota + 1
	QQ
	Alipay
)

// ExtractBrowser 从UA提取浏览器
func ExtractBrowser(useragent string) Browser {
	var browser Browser
	useragent = strings.ToLower(useragent)
	switch {
	case strings.Contains(useragent, "micromessenger"):
		browser = Wechat
	case strings.Contains(useragent, "qq"):
		browser = QQ
	case strings.Contains(useragent, "alipay"):
		browser = Alipay
	}
	return browser
}

// GetBrowserFromCtx 获取浏览器标识
func GetBrowserFromCtx(ctx context.Context) Browser {
	if u, ok := ctx.Value(CtxUaKey).(UserAgent); ok {
		return u.Browser
	}
	return Browser(0)
}

func (b Browser) IsWechat() bool {
	return b == Wechat
}

func (b Browser) IsQQ() bool {
	return b == QQ
}

func (b Browser) IsAlipay() bool {
	return b == Alipay
}

func (b Browser) IsOther() bool {
	return !b.IsWechat() && !b.IsQQ() && !b.IsAlipay()
}
