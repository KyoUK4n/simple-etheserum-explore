package limiter

import (
	"sync"
	"time"
)

type SlidingWindowLimiter struct {
	windowSize  time.Duration
	lock        sync.Mutex
	maxRequests int
	timestamps  []time.Time
}

const (
	DefaultRequestPeriod = time.Second
	DefaultRequestLimit  = 5
)

func NewSlidingWindowLimiter(windowSize time.Duration, maxRequests int) *SlidingWindowLimiter {
	if windowSize <= 0 {
		windowSize = DefaultRequestPeriod
	}
	if maxRequests <= 0 {
		maxRequests = DefaultRequestLimit
	}
	return &SlidingWindowLimiter{
		windowSize:  windowSize,
		lock:        sync.Mutex{},
		maxRequests: maxRequests,
		timestamps:  make([]time.Time, 0),
	}
}

func (s *SlidingWindowLimiter) Allow() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	now := time.Now()
	windowStart := now.Add(-s.windowSize)

	s.evict(windowStart)

	if len(s.timestamps) >= s.maxRequests {
		return false
	}
	s.timestamps = append(s.timestamps, now)
	return true
}

func (s *SlidingWindowLimiter) evict(windowStart time.Time) {
	i := 0
	for i < len(s.timestamps) && s.timestamps[i].Before(windowStart) {
		i++
	}
	s.timestamps = s.timestamps[i:]
}

// Count 返回当前窗口的请求数
func (s *SlidingWindowLimiter) Count() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	windowStart := time.Now().Add(-s.windowSize)
	s.evict(windowStart)
	return len(s.timestamps)
}

// Remaining 返回当前窗口内剩余可用请求数
func (s *SlidingWindowLimiter) Remaining() int {
	remaining := s.maxRequests - s.Count()
	if remaining < 0 {
		return 0
	}
	return remaining
}
