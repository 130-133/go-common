package appleiap

import "testing"

func TestNewApi(t *testing.T) {
	a := NewApi(StoreAuthConfig{
		Bid:        "com.minitech.miniaixue",
		Iss:        "090b0353-f50e-4cf9-aa4f-225aa94020b0",
		Kid:        "P7KVC29TPJ",
		Kp8:        "D:/work/sirius-unionpay/cert/SubscriptionKey_P7KVC29TPJ.p8",
		PublicCert: "D:/work/sirius-unionpay/cert/AppleRootCA-G3.pem",
	})
	a.LookUp("MVWMQ30VW4")
}
