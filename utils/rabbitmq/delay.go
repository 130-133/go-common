package rabbitmq

type DelayTime int64

var (
	Delay1s  DelayTime = 1000            // 1秒
	Delay3s            = 3 * Delay1s     // 3秒
	Delay10s           = 10 * Delay1s    // 10秒
	Delay30s           = 30 * Delay1s    // 30秒
	Delay1m            = 60 * Delay1s    // 1分钟
	Delay5m            = 300 * Delay1s   // 5分钟
	Delay10m           = 600 * Delay1s   // 10分钟
	Delay30m           = 1800 * Delay1s  // 30分钟
	Delay1h            = 3600 * Delay1s  // 1小时
	Delay3h            = 10800 * Delay1s // 3小时
	Delay6h            = 21600 * Delay1s // 6小时
	Delay12h           = 43200 * Delay1s // 12小时
	Delay1d            = 86400 * Delay1s // 1天
)

//延迟阶梯
var delayStage = []DelayTime{
	Delay1s,
	Delay3s,
	Delay10s,
	Delay30s,
	Delay1m,
	Delay5m,
	Delay10m,
	Delay30m,
	Delay1h,
	Delay3h,
	Delay6h,
	Delay12h,
	Delay1d,
}

var (
	delayDlX             = "sir@go_delay_dlx"   //死信交换机
	delayQueue           = `sir@go_delay_%ds`   //延迟队列
	delayExchangeDefault = `sir@go_delay_ex_%d` //延迟队列默认交换机
	nextDelayTimeKey     = "next_delay_time"    //header key 标记下阶段重试时长
)

// Next 返回当前延迟时间的下个阶段
func (d DelayTime) Next() DelayTime {
	max := len(delayStage) - 1
	for k := 0; k < max; k++ {
		if delayStage[k] >= d && d < delayStage[k+1] {
			return delayStage[k+1]
		}
	}
	return -1
}
func (d DelayTime) Int64() int64 {
	return int64(d)
}
