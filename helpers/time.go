package helpers

import (
	"time"
)

// Timer 计时器数据结构
type Timer struct {
	duration time.Duration
	expireAt int64
}

// NewTimer 使用指定的时间间隔新建一个计时器
func NewTimer(duration time.Duration) *Timer {
	expireAt := int64(0)
	if duration < 0 {
		expireAt = -1
	}
	return &Timer{duration, expireAt}
}

// IsTimeUp 计时器是否已超时
func (t *Timer) IsTimeUp() bool {
	return t.now() >= t.expireAt
}

// Renew 使用创建时指定的时间间隔重新设置计时器
func (t *Timer) Renew() {
	t.expireAt = getTimestamp(time.Now().Add(t.duration))
}

// IsEnabled 计时器当前是否生效
func (t *Timer) IsEnabled() bool {
	return t.expireAt >= 0
}

func (t *Timer) now() int64 {
	return getTimestamp(time.Now())
}

func getTimestamp(t time.Time) int64 {
	return t.UTC().Unix()
}
