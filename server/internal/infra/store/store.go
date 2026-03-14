package store

import (
	"context"
	"sync"

	"github.com/AntonioMartinezLopez/enginsight/server/internal/domain/counter"
)

type CounterStore struct {
	mutex sync.RWMutex
	count counter.Count
}

func New() *CounterStore {
	return &CounterStore{}
}

func (s *CounterStore) Increment(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count++
	return nil
}

func (s *CounterStore) GetMessageCount(ctx context.Context) (counter.Count, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.count, nil
}
