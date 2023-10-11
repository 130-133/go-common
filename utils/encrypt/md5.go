package encrypt

import (
	"crypto/md5"
	"fmt"
	"strings"
)

type IMD5 interface {
	Encode(data string, options ...CommonOpt) string
	EncodeToken(data string, options ...CommonOpt) string
	EncodeSalt(data string, salts ...string) string
}

type MD5Encrypt struct {
	opt   CommonOptions
	token string
}

func NewMD5(tokens ...string) IMD5 {
	token := strings.Join(tokens, "")
	return &MD5Encrypt{
		token: token,
	}
}

func (e *MD5Encrypt) Encode(data string, options ...CommonOpt) string {
	return e.md5(data, options...)
}

func (e *MD5Encrypt) EncodeToken(data string, options ...CommonOpt) string {
	text := fmt.Sprintf("%s%s", data, e.token)
	return e.md5(text, options...)
}

func (e *MD5Encrypt) EncodeSalt(data string, salts ...string) string {
	salt := strings.Join(salts, "")
	text := fmt.Sprintf("%s%s", data, salt)
	return e.md5(text)
}

func (e *MD5Encrypt) md5(str string, options ...CommonOpt) string {
	for _, o := range options {
		o.Apply(&e.opt)
	}
	newSig := md5.Sum([]byte(str))
	newArr := fmt.Sprintf("%x", newSig)
	return turnCase(newArr, e.opt.signCase)
}
