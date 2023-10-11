package help

import (
	"testing"
)

func TestNewSet(t *testing.T) {
	aa := []int64{
		1, 3, 5, 768, 3, 31, 1,
	}
	a := NewSet(aa...)
	t.Logf("%#v\n", a.Values())
	a.Add(2).Add(3)
	t.Logf("%#v\n", a.Values())

	var b []int64
	b = a.Values()
	func(bb []int64) {
		for _, x := range bb {
			t.Logf("%d\n", x)
		}
	}(b)
}

func TestSet_Add(t *testing.T) {
	a := NewSet[int64]()
	a.Add(2).Add(3)
	t.Logf("%#v\n", a.Values())
}
