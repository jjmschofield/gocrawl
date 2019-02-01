package counters

import (
	"sync/atomic"
)

type AtomicInt64 struct {
	count int64
	peak  int64
}

func (counter *AtomicInt64) Add(delta int64) int64 {
	count := atomic.AddInt64(&counter.count, delta)

	if count > atomic.LoadInt64(&counter.peak) {
		atomic.SwapInt64(&counter.peak, count)
	}

	return count
}

func (counter *AtomicInt64) Sub(delta int64) int64 {
	return atomic.AddInt64(&counter.count, -delta)
}

func (counter *AtomicInt64) Count() int64 {
	return atomic.LoadInt64(&counter.count)
}

func (counter *AtomicInt64) Peak() int64 {
	return atomic.LoadInt64(&counter.peak)
}
