package logger

import (
	"errors"
	"io"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type MRotate struct {
	PathStr      string
	PathNameStr  string
	RotationTime time.Duration
	KeepDays     time.Duration
}

func (l *MRotate) New() (io.Writer, error) {
	if l.PathStr == "" || l.PathNameStr == "" {
		return nil, errors.New("no path")
	}

	fileName := path.Join(l.PathStr, l.PathNameStr)
	return rotatelogs.New(
		fileName+".%Y%m%d",
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间(天)
		rotatelogs.WithMaxAge(l.KeepDays),
		// 设置日志切割时间间隔(天)
		rotatelogs.WithRotationTime(l.RotationTime),
	)
}
