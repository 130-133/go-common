package encrypt

import "testing"

func TestNewJWT(t *testing.T) {
	token := "my_secret_key"
	jwt := NewJwt(JwtConfig{Token: token})
	t.Log(jwt)
}

func TestJWTEncrypt_Encode(t *testing.T) {
	token := "my_secret_key"
	jwt := NewJwt(JwtConfig{Token: token})
	waitData := map[string]interface{}{
		"my":     "test",
		"number": 1,
	}
	encode := jwt.Encode(waitData)
	t.Log(encode.String())
	t.Log(encode.Error())
}

func TestJWTEncrypt_Decode(t *testing.T) {
	token := "my_secret_key"
	jwt := NewJwt(JwtConfig{Token: token})
	waitData := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEsIm15IjoidGVzdCIsIm51bWJlciI6MX0.e9dHz1jZfn6P2d5vODn8fXevDUt5gDkip3klvmbhSm4"
	decode := jwt.Decode(waitData)
	t.Log(decode.Verify())
	t.Log(decode.Data())
	t.Log(decode.Error())
}
