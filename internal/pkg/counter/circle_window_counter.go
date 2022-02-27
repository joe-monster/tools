// @Title  circle_window_counter
// @Description 滑动窗口计数器，暂时只计数，未增加限制相关逻辑
// @Author Joe
// @Create 2022-02-26
// @Update 2022-02-26
package counter

import (
	"errors"
	"math"
	"sync"
	"time"
)

type circleWindowCounter struct {
	w *sync.Mutex
	windowSize int64	//窗口大小，单位：毫秒
	bucketNum int64 //小窗口数量
	bucketSize int64	//小窗口大小，这个数值由窗口大小 / 小窗口切分数量决定，不允许出现除不尽的现象
	counters []int	//小窗口计数器，初始化时，根据小窗口切分数量来确定数组大小和容量
	index int64	//当前小窗口计数器的索引
	startTimestampNano int64	//窗口开始时间
}

//bucket表示小窗口切分数量，这个数值决定了平滑度
func NewCircleWindow(size, bucket int64) (*circleWindowCounter, error) {

	if t := size % bucket;t != 0 {
		return nil, errors.New("窗口大小必须可以被切分数量除尽")
	}

	var c = circleWindowCounter{
		w: new(sync.Mutex),
		windowSize: size,
		bucketNum: bucket,
		bucketSize: size / bucket,
		counters: make([]int, bucket, bucket),
		index: 0,
		startTimestampNano: time.Now().UnixNano() / 1e6 - size,	//哨兵，提前排空窗口内情况，收到流量时可以确保滑动窗口数肯定大于0
	}

	return &c, nil
}

//该方法必须并发安全
func (c *circleWindowCounter) Set() int {
	now := time.Now().UnixNano() / 1e6

	c.w.Lock()
	defer c.w.Unlock()

	//计算需要滑动几个窗口
	slideBucketNum := (now - c.startTimestampNano - c.windowSize) / c.bucketSize
	if slideBucketNum > 0 {
		maxSlideNum := int(math.Min(float64(slideBucketNum), float64(c.bucketNum)))	//循环的主要工作为清除相关位置的计数器，计数器最多不会超过bucketNum，这里做一个
		//因为有哨兵的作用，滑动的窗口数肯定会大于0，触发滑动窗口
		for i := 0; i < maxSlideNum; i++ {	//超多少循环多少次
			c.index = (c.index + 1) % c.bucketNum
			c.counters[c.index] = 0
		}
		c.startTimestampNano += slideBucketNum * c.bucketSize	//滑动了几个小窗口，就把起算时间往后推几个小窗口
	}

	c.counters[c.index]++

	return c.counters[c.index]
}
