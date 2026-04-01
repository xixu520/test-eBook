package worker

import (
	"math"
	"time"
)

// RetryPolicy 指数退避重试策略
type RetryPolicy struct {
	MaxRetries    int           // 最大重试次数，默认 5
	InitialDelay  time.Duration // 首次延迟，默认 5s
	MaxDelay      time.Duration // 最大延迟，默认 5min
	BackoffFactor float64       // 退避系数，默认 2.0
}

// DefaultRetryPolicy 返回默认的重试策略
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxRetries:    5,
		InitialDelay:  5 * time.Second,
		MaxDelay:      5 * time.Minute,
		BackoffFactor: 2.0,
	}
}

// NewRetryPolicy 从配置创建重试策略
func NewRetryPolicy(maxRetries int, initialDelaySec int) RetryPolicy {
	p := DefaultRetryPolicy()
	if maxRetries > 0 {
		p.MaxRetries = maxRetries
	}
	if initialDelaySec > 0 {
		p.InitialDelay = time.Duration(initialDelaySec) * time.Second
	}
	return p
}

// NextDelay 根据当前重试次数计算下次延迟
// 公式: min(InitialDelay * BackoffFactor^retryCount, MaxDelay)
// 示例: 5s → 10s → 20s → 40s → 80s (capped at 5min)
func (p *RetryPolicy) NextDelay(retryCount int) time.Duration {
	delay := float64(p.InitialDelay) * math.Pow(p.BackoffFactor, float64(retryCount))
	if delay > float64(p.MaxDelay) {
		return p.MaxDelay
	}
	return time.Duration(delay)
}

// ShouldRetry 是否应该继续重试
func (p *RetryPolicy) ShouldRetry(retryCount int) bool {
	return retryCount < p.MaxRetries
}

// NextRetryTime 计算下次重试的绝对时间
func (p *RetryPolicy) NextRetryTime(retryCount int) time.Time {
	return time.Now().Add(p.NextDelay(retryCount))
}
