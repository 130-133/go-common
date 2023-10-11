package help

import "testing"

func TestPrice(t *testing.T) {
	t.Logf("%#v\n", CentToYuan(12345).Round(2))
	t.Logf("%#v\n", CentToYuan(12300).Str())
	t.Logf("%#v\n", CentToYuan(12345).Float64())
	t.Logf("%#v\n", CentToYuan(12345).Int())
	t.Logf("%#v\n", CentToYuan(1234500).Float64())
	t.Logf("%#v\n", CentToYuan(1234500).Round(2))
	t.Logf("%#v\n", CentToYuan(0).Str())
	t.Logf("%#v\n", CentToYuan(0).StrOrEmpty())
	t.Logf("%#v\n", CentToYuan(0).Int())
}

func TestCentToYuan(t *testing.T) {
	t.Log(CentToYuan(100))
	t.Log(CentToYuan(112354789.4123))
}

func TestYuanToCent(t *testing.T) {
	t.Log(YuanToCent(1))
	t.Log(YuanToCent(100))
	t.Log(YuanToCent(100.10))
	t.Logf("%#v\n", YuanToCent(100.1001))
	t.Logf("%#v\n", YuanToCent(100.1001).Round(3))
	t.Logf("%#v\n", YuanToCent(100.1001).Int())
	t.Logf("%#v\n", YuanToCent(100.1001).Str())
}

func TestToPrice(t *testing.T) {
	t.Logf("%#v\n", ToPrice("-123.4532").Str())
	t.Logf("%#v\n", ToPrice("-123.4532").Round(3))
	t.Logf("%#v\n", ToPrice("123.4532").Int())
	t.Logf("%#v\n", ToPrice(" 123.4532 ").Float64())
	t.Logf("%#v\n", ToPrice("23,123.4532"))
}
