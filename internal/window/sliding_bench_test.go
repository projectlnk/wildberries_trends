package window

import (
	"testing"
	"time"
)

func BenchmarkSlidingWindow_Add(b *testing.B) {
	win := NewSlidingWindow(5*time.Minute, 5*time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		win.Add("test")
	}
}

func BenchmarkSlidingWindow_GetAllCounts(b *testing.B) {
	win := NewSlidingWindow(5*time.Minute, 5*time.Second)
	for i := 0; i < 1000; i++ {
		win.Add("query")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		win.GetAllCounts()
	}
}
