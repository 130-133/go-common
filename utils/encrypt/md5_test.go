package encrypt

import "testing"

func TestMD5Encrypt_Encode(t *testing.T) {
	t.Log(NewMD5().Encode("abc"))
	t.Log(NewMD5().Encode("abc", SetUpper()))
	t.Log(NewMD5().Encode("abc", SetLower()))
}
