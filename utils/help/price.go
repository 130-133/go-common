package help

import (
	"fmt"
	"strconv"
	"strings"
)

type PriceType interface {
	~float64 | ~float32 | ~int | ~int64 | ~int32
}
type Price float64

// CentToYuan 分转元
func CentToYuan[T PriceType](cent T) Price {
	return Price(float64(cent) / 100)
}

// YuanToCent 元转分
func YuanToCent[T PriceType](yuan T) Price {
	return Price(float64(yuan) * 100)
}

// ToPrice 字符串转金额
func ToPrice(str string) Price {
	str = strings.ReplaceAll(str, ",", "")
	str = strings.TrimSpace(str)
	i, _ := strconv.ParseFloat(str, 10)
	return Price(i)
}

func (p Price) Round(num int) string {
	format := fmt.Sprintf("%%.%df", num)
	return fmt.Sprintf(format, p)
}

func (p Price) Str() string {
	str := fmt.Sprintf("%f", p)
	if strings.Contains(str, ".") {
		return strings.TrimSuffix(strings.TrimRight(str, "0"), ".")
	}
	return str
}
func (p Price) Int() int {
	return int(p)
}
func (p Price) Int64() int64 {
	return int64(p)
}
func (p Price) Int32() int32 {
	return int32(p)
}
func (p Price) Float64() float64 {
	return float64(p)
}
func (p Price) Float32() float32 {
	return float32(p)
}

// StrOrEmpty 0时返回空值 反之输出字符串
func (p Price) StrOrEmpty() string {
	if p == 0 {
		return ""
	}
	return p.Str()
}

func (p Price) CentToYuan() Price {
	return p / 100
}
func (p Price) YuanToCent() Price {
	return p * 100
}
