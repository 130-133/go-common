package help

import (
	"testing"
	"time"
)

func TestGoWait(t *testing.T) {
	var fn []func()
	for i := int64(1e5); i > 0; i-- {
		a := func() {
			t.Log("start")
			time.Sleep(time.Duration(3) * time.Second)
			t.Log("end")
		}
		fn = append(fn, a)
	}
	GoWait(fn...)
}

func TestGoForLimit(t *testing.T) {
	a := []string{
		"1", "2", "5", "7", "8",
		"1", "2", "5", "7", "8",
		"1", "2", "5", "7", "8",
		"1", "2", "5", "7", "8",
		"1", "2", "5", "7", "8",
		"1", "2", "5", "7", "8",
	}
	GoForeachLimit(a, 5, func(i int) {
		time.Sleep(3 * time.Second)
		t.Log(a[i])
	})
}
