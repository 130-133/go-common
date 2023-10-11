package encrypt

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

type ISha256 interface {
	Encode(data string, options ...CommonOpt) string
	EncodeToken(data string, options ...CommonOpt) string
	EncodeSalt(data string, salts ...string) string
}

type Sha256Encrypt struct {
	opt   CommonOptions
	token string
}

func NewSha256(tokens ...string) ISha256 {
	token := strings.Join(tokens, "")
	return &Sha256Encrypt{
		token: token,
	}
}

func (e *Sha256Encrypt) Encode(data string, options ...CommonOpt) string {
	return e.sha256(data, options...)
}

func (e *Sha256Encrypt) EncodeToken(data string, options ...CommonOpt) string {
	text := fmt.Sprintf("%s%s", data, e.token)
	return e.sha256(text, options...)
}

func (e *Sha256Encrypt) EncodeSalt(data string, salts ...string) string {
	salt := strings.Join(salts, "")
	text := fmt.Sprintf("%s%s", data, salt)
	return e.sha256(text)
}

func (e *Sha256Encrypt) sha256(str string, options ...CommonOpt) string {
	for _, o := range options {
		o.Apply(&e.opt)
	}
	newSig := sha256.New().Sum([]byte(str))
	newArr := fmt.Sprintf("%x", newSig)
	return turnCase(newArr, e.opt.signCase)
}
