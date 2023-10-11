package encrypt

import "testing"

func TestSha256Encrypt_Encode(t *testing.T) {
	t.Log(NewSha256().Encode("abc"))
	t.Log(NewSha256().Encode("abc", SetUpper()))
	t.Log(NewSha256().Encode("abc", SetLower()))
}
