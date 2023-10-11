package help

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func IntAbs(n int64) int64 {
	if n >= 0 {
		return n
	}
	return int64(math.Abs(float64(n)))
}

func GetMethodName() string {
	if pc, _, _, ok := runtime.Caller(0); ok {
		f := runtime.FuncForPC(pc)
		return f.Name()
	}
	return ""
}

// TimeFormatStrict 时间格式 严格模式
func TimeFormatStrict(ts time.Time) string {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	todayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterdayZero := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	thisYearZero := time.Date(now.Year(), 0, 0, 0, 0, 0, 0, now.Location())
	if sub := ts.Sub(todayZero).Seconds(); sub >= 0 && sub < 86400 {
		return ts.Format("今天 15:04")
	} else if sub := ts.Sub(yesterdayZero).Seconds(); sub >= 0 && sub < 172800 {
		return ts.Format("昨天 15:04")
	} else if ts.Unix() > thisYearZero.Unix() {
		return ts.Format("01-02 15:04")
	} else {
		return ts.Format("2006-01-02 15:04")
	}
}

func CountDisplay(count int64, min int64) (result string) {
	if count > min && count < 10000 {
		result = strconv.FormatInt(count, 10)
	} else if count >= 10000 && count < 1000000 {
		result = fmt.Sprintf("%.1f", float64(count)/10000) + "万"
	} else if count >= 1000000 && count < 100000000 {
		result = strconv.FormatInt(count/10000, 10) + "万"
	} else if count >= 100000000 {
		result = strconv.FormatInt(count/100000000, 10) + "亿"
	}
	return
}

// FetchVersionNum 3.10.1 >>> 3010001
func FetchVersionNum(version string) (n int64) {
	for i, v := range strings.Split(version, ".") {
		num, _ := strconv.ParseInt(v, 10, 64)
		if i < 2 {
			n += num * int64(math.Pow(10, float64(6/(i+1))))
		} else {
			n += num
		}
	}
	return
}
