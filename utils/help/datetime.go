package help

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

const (
	FormatTime        = "2006-01-02 15:04:05"
	FormatDate        = "2006-01-02"
	FormatCnDate      = "2006年01月02日"
	FormatCnYear      = "2006年"
	FormatCnYearMonth = "2006年01月"
	FormatRawTime     = "20060102150405"
	FormatRawMonth    = "200601"
	FormatRawDate     = "20060102"
)

func DateFormat(arg interface{}, format string) string {
	t := reflect.TypeOf(arg)
	var last time.Time
	switch t.Name() {
	case "NullTime":
		data := arg.(sql.NullTime)
		if data.Valid == true {
			last = data.Time
		}
	case "Time":
		last = arg.(time.Time)
	}
	if last.IsZero() {
		return ""
	}
	return last.Local().Format(format)
}

var tryFormatDate = []string{
	time.RFC3339,
	FormatTime,
	FormatRawTime,
	FormatDate,
	FormatRawDate,
}

// ParseDate 解析字符串时间格式
func ParseDate(date string) time.Time {
	for _, format := range tryFormatDate {
		if timeData, err := time.ParseInLocation(format, date, time.Local); err == nil {
			return timeData
		}
	}
	return time.Time{}
}

// ParseUnix 解析时间戳
func ParseUnix(unix int64) time.Time {
	if unix == 0 {
		return time.Time{}
	}
	unixStr := fmt.Sprintf("%d", unix)
	switch len(unixStr) {
	case 10:
		return time.Unix(unix, 0)
	case 13:
		return time.UnixMilli(unix)
	case 16:
		return time.UnixMicro(unix)
	}
	return time.Time{}
}

// ParseUnixWithNullTime 解析时间戳 0时间返回null
func ParseUnixWithNullTime(unix int64) sql.NullTime {
	times := ParseUnix(unix)
	if times.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: times, Valid: true}
}

// NullTimeToUnix 可null时间类型转时间戳
func NullTimeToUnix(nullTime sql.NullTime) int64 {
	if !nullTime.Valid || nullTime.Time.IsZero() {
		return 0
	}
	return nullTime.Time.Unix()
}
