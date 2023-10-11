package help

import "testing"

func TestToAny(t *testing.T) {
	t.Logf("%+v", ToAny("321"))
	t.Logf("%+v", ToAny([]byte("1234")))
}

func TestStrToAny(t *testing.T) {
	a := ToAny("1333333333333333333333333333333333333")
	t.Log(a.Int())
	t.Log(a.Int32())
	t.Log(a.Int64())
	t.Log(a.Float32())
	t.Log(a.Float64())
	t.Log(a.Bool())
	t.Log(a.Slice())
}

func TestToString(t *testing.T) {
	t.Log(ToString(nil))
	t.Log(ToString("this is a test"))
	t.Log(ToString(1))
	t.Log(ToString([]uint8{1, 3, 4}))
	t.Log(ToString([]byte("this is a test")))
	t.Log(ToString([2]int{1, 2}))
}
