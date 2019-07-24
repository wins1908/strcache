package strcache

import (
	"context"
	"sync"
)

func NewMemoryFetcher() *MemoryFetcher {
	return &MemoryFetcher{values: make(map[string]string)}
}

var _ Fetcher = (*MemoryFetcher)(nil)

type MemoryFetcher struct {
	mu     sync.Mutex
	values map[string]string
}

func (s *MemoryFetcher) Fetch(ctx context.Context, key string, notFoundFn NotFoundFunc, foundFn FoundFunc) error {
	s.mu.Lock()
	if value, exists := s.values[key]; exists {
		s.mu.Unlock()
		return foundFn(ctx, value, false)
	}

	newFnCalled := false
	err := notFoundFn(ctx, NewValueFunc(func(ctx context.Context, newValue string) error {
		s.values[key] = newValue
		s.mu.Unlock()
		newFnCalled = true
		return foundFn(ctx, newValue, true)
	}))

	if !newFnCalled {
		s.mu.Unlock()
	}

	return err
}

func (s *MemoryFetcher) Clear(ctx context.Context, key string) error {
	s.mu.Lock()
	if _, exists := s.values[key]; exists {
		delete(s.values, key)
	}
	s.mu.Unlock()
	return nil
}
