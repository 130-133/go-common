package encrypt

import (
	"testing"
	"time"
)

func TestParamSign(t *testing.T) {
	sign := NewParamSign("522e6ccde68f4584a3282ce0b0eca1cf")
	a := sign.Sign(map[string]interface{}{
		"a": 1,
		"b": 2,
	})
	t.Log(a.SignedVal())
}

func TestParamSign_Format(t *testing.T) {
	sign := NewParamSign("522e6ccde68f4584a3282ce0b0eca1cf")
	sign.SetFormatter(func(k, v string) string {
		return v
	})
	a := sign.Sign(map[string]interface{}{
		"a": 1,
		"b": 2,
	})
	t.Log(a.SignedVal())
}

func TestParamSignEncrypt_Case(t *testing.T) {
	sign := NewParamSign("522e6ccde68f4584a3282ce0b0eca1cf")
	a := sign.SetUpperSign().Sign(map[string]interface{}{
		"a":        1,
		"b":        2,
		"signtime": time.Now().Unix(),
	})
	t.Log(a.SignedVal())
}

func TestParamSignResult_Verify(t *testing.T) {
	sign := NewParamSign("522e6ccde68f4584a3282ce0b0eca1cf")
	a := sign.Sign(map[string]interface{}{
		"a":        1,
		"b":        2,
		"signtime": time.Now().Unix(),
	})
	t.Log(a.SignedVal())
	t.Log(a.Verify())
}
