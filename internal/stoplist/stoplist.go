package stoplist

import (
	"sync"
)

type StopList struct {
	mu    sync.RWMutex
	words map[string]struct{}
}

func New() *StopList {
	return &StopList{
		words: make(map[string]struct{}),
	}
}

// Add добавляет слово в стоп-лист (регистронезависимо, можно привести к нижнему регистру)
func (s *StopList) Add(word string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.words[word] = struct{}{}
}

// Remove удаляет слово из стоп-листа
func (s *StopList) Remove(word string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.words, word)
}

// Contains проверяет, есть ли слово в стоп-листе
func (s *StopList) Contains(word string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.words[word]
	return ok
}

// List возвращает все слова из стоп-листа (для отладки)
func (s *StopList) List() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]string, 0, len(s.words))
	for w := range s.words {
		res = append(res, w)
	}
	return res
}
