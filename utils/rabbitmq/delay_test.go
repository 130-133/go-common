package rabbitmq

import "testing"

func TestDelayTime_Next(t *testing.T) {
	t.Log(Delay1s.Next())
	t.Log(Delay3s.Next())
	t.Log(Delay10s.Next())
	t.Log(Delay30m.Next())
	t.Log(Delay1h.Next())
	t.Log(Delay1d.Next())
	t.Log(DelayTime(23000).Next())
}
