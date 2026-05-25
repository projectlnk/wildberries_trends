package top

import (
	"container/heap"
)

// Pair хранит запрос и его частоту
type Pair struct {
	Query string
	Count int64
}

// MinHeap реализует интерфейс heap для хранения top N (наименьшие в корне)
type MinHeap []Pair

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].Count < h[j].Count }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(Pair))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// GetTopN возвращает топ-N запросов из мапы частот
// n - количество запрашиваемых элементов (если больше уникальных – вернёт все)
func GetTopN(counts map[string]int64, n int) []string {
	if n <= 0 {
		return []string{}
	}
	h := &MinHeap{}
	heap.Init(h)

	for query, cnt := range counts {
		if h.Len() < n {
			heap.Push(h, Pair{query, cnt})
		} else if cnt > (*h)[0].Count {
			heap.Pop(h)
			heap.Push(h, Pair{query, cnt})
		}
	}
	// Извлекаем из кучи в порядке возрастания частот, затем разворачиваем
	result := make([]string, h.Len())
	for i := h.Len() - 1; i >= 0; i-- {
		result[i] = heap.Pop(h).(Pair).Query
	}
	return result
}
