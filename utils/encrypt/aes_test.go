package encrypt

import (
	"encoding/json"
	"testing"
)

func TestAesEncrypt_Encode(t *testing.T) {
	a := NewAes("adFEdE2A10SO2022", "1234567891234567", OutHex)
	b, _ := json.Marshal(map[string]interface{}{
		"a": 1,
	})
	c := a.OutType(OutHex).Encode(b)
	t.Log(c.String())
	t.Log(c.Error())
}

func TestAesEncrypt_Decode(t *testing.T) {
	a := NewAes("adFEdE2A10SO2022", "1234567891234567", OutHex)
	b := a.OutType(OutHex).Decode("cf382cb3e4f6b03d40118f172d3bf3f1")
	t.Log(b.Data())
	t.Log(b.Error())
}
