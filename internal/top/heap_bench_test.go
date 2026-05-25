package top

import (
	"testing"
)

func BenchmarkGetTopN(b *testing.B) {
	// Генерируем 1000 уникальных запросов со случайными частотами
	counts := make(map[string]int64)
	for i := 0; i < 1000; i++ {
		counts[string(rune(i))] = int64(i % 100)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetTopN(counts, 10)
	}
}
