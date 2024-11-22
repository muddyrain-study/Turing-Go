package time

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	now := time.Now()
	t.Log(now)

	year := now.Year()     //年
	month := now.Month()   //月
	day := now.Day()       //日
	hour := now.Hour()     //小时
	minute := now.Minute() //分钟
	second := now.Second() //秒
	//02d输出的整数不足两位 用0补足
	fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", year, month, day, hour, minute, second)
}

func TestTimeUnix(t *testing.T) {
	now := time.Now()
	t.Log(now)

	timestamp := now.Unix()         //时间戳
	milli := now.UnixMilli()        //毫秒时间戳
	micro := now.UnixMicro()        //微秒时间戳
	timestampNano := now.UnixNano() //纳秒时间戳
	t.Log(timestamp, milli, micro, timestampNano)
}

func TestTimeParse(t *testing.T) {
	timeStr := "2021-07-26 12:55:44"
	timeObj, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		t.Error(err)
	}
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Error(err)
	}
	timeObj = timeObj.In(location)
	t.Log(timeObj)
}

func TestTimeFormat(t *testing.T) {
	now := time.Now()
	t.Log(now)

	timeStr := now.Format("2006年01月02日 15时04分05秒") // 24小时制
	t.Log(timeStr)
}
func TestTimeDuration(t *testing.T) {
	fiveMinute := 5 * time.Minute

	now := time.Now()

	later := now.Add(fiveMinute)
	t.Log(later)

	sub := later.Sub(now)
	t.Log(sub)
}
func TestEqualBeforeAfter(t *testing.T) {
	now := time.Now()
	later := now.Add(5 * time.Minute)
	t.Log(now.Equal(later))
	t.Log(now == later)
	t.Log(now == now)
	t.Log(now.Before(later))
	t.Log(now.After(later))
}

func TestTicker(t *testing.T) {
	//	tick := time.Tick(5 * time.Second)
	//	fmt.Println(time.Now())
	//	for i := range tick {
	//		t.Log(i)
	//	}
	//fmt.Println(time.Now())
	//time.AfterFunc(5*time.Second, func() {
	//	fmt.Println("5s后执行")
	//})
	//for {
	//
	//}
	var wg sync.WaitGroup
	wg.Add(2)
	//NewTimer 创建一个 Timer，它会在最少过去时间段 d 后到期，向其自身的 C 字段发送当时的时间
	timer1 := time.NewTimer(2 * time.Second)

	//NewTicker 返回一个新的 Ticker，该 Ticker 包含一个通道字段，并会每隔时间段 d 就向该通道发送当时的时间。它会调
	//整时间间隔或者丢弃 tick 信息以适应反应慢的接收者。如果d <= 0会触发panic。关闭该 Ticker 可
	//以释放相关资源。
	ticker1 := time.NewTicker(2 * time.Second)

	go func(t *time.Ticker) {
		defer wg.Done()
		for {
			<-t.C
			fmt.Println("get ticker1", time.Now().Format("2006-01-02 15:04:05"))
		}
	}(ticker1)

	go func(t *time.Timer) {
		defer wg.Done()
		for {
			<-t.C
			fmt.Println("get timer", time.Now().Format("2006-01-02 15:04:05"))
			//Reset 使 t 重新开始计时，（本方法返回后再）等待时间段 d 过去后到期。如果调用时t
			//还在等待中会返回真；如果 t已经到期或者被停止了会返回假。
			//t.Reset(2 * time.Second)
		}
	}(timer1)

	wg.Wait()
}
