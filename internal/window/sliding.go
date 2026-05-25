package window

import (
	"sync"
	"time"

	"wildberries-trends/internal/metrics"
)

type SlidingWindow struct {
	mu           sync.Mutex
	slots        []map[string]int64
	slotDuration time.Duration
	windowSize   int
	currentSlot  int
	lastSlotTime time.Time
}

func NewSlidingWindow(period, slotDuration time.Duration) *SlidingWindow {
	size := int(period / slotDuration)
	if size < 1 {
		size = 1
	}
	slots := make([]map[string]int64, size)
	for i := 0; i < size; i++ {
		slots[i] = make(map[string]int64)
	}
	return &SlidingWindow{
		slots:        slots,
		slotDuration: slotDuration,
		windowSize:   size,
		currentSlot:  0,
		lastSlotTime: time.Now(),
	}
}

func (sw *SlidingWindow) rotate(now time.Time) {
	elapsed := now.Sub(sw.lastSlotTime)
	shift := int(elapsed / sw.slotDuration)
	if shift == 0 {
		return
	}
	if shift >= sw.windowSize {
		for i := 0; i < sw.windowSize; i++ {
			sw.slots[i] = make(map[string]int64)
		}
		sw.currentSlot = 0
		sw.lastSlotTime = now
		return
	}
	for i := 1; i <= shift; i++ {
		idx := (sw.currentSlot + i) % sw.windowSize
		sw.slots[idx] = make(map[string]int64)
	}
	sw.currentSlot = (sw.currentSlot + shift) % sw.windowSize
	sw.lastSlotTime = sw.lastSlotTime.Add(time.Duration(shift) * sw.slotDuration)
}

func (sw *SlidingWindow) Add(query string) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	now := time.Now()
	sw.rotate(now)
	sw.slots[sw.currentSlot][query]++
}

func (sw *SlidingWindow) GetAllCounts() map[string]int64 {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	now := time.Now()
	sw.rotate(now)
	res := make(map[string]int64)
	for _, slot := range sw.slots {
		for q, c := range slot {
			res[q] += c
		}
	}
	// Обновляем метрику уникальных запросов в окне
	metrics.UniqueQueriesInWindow.Set(float64(len(res)))
	return res
}
