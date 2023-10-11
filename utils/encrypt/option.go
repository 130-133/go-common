package encrypt

import "strings"

type SignCase int

const (
	Lower SignCase = iota + 1
	Upper
)

type CommonOptions struct {
	signCase SignCase
}

type CommonOpt func(*CommonOptions)

func (o CommonOpt) Apply(options *CommonOptions) {
	o(options)
}

func SetLower() CommonOpt {
	return func(options *CommonOptions) {
		options.signCase = Lower
	}
}

func SetUpper() CommonOpt {
	return func(options *CommonOptions) {
		options.signCase = Upper
	}
}

func turnCase(str string, signCase SignCase) string {
	switch signCase {
	case Lower:
		return strings.ToLower(str)
	case Upper:
		return strings.ToUpper(str)
	default:
		return str
	}
}
