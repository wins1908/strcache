package strcache

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const redisCacheKey = "cache:{key}"

var _ Fetcher = (*RedisFetcher)(nil)

type RedisFetcher struct {
	mu         sync.Mutex
	redis      *redis.Client
	expiration time.Duration
}

func (s *RedisFetcher) Fetch(
	ctx context.Context,
	key string,
	notFoundFn NotFoundFunc,
	foundFn FoundFunc,
) error {
	var err error
	var value string
	s.mu.Lock()
	cacheKey := s.cacheKey(key)
	c := s.redis.WithContext(ctx)
	value, err = c.Get(cacheKey).Result()
	if err == nil {
		s.mu.Unlock()
		return foundFn(ctx, value, true)
	}

	if err != redis.Nil {
		logrus.Error(err)
	}

	newFnCalled := false
	err = notFoundFn(ctx, NewValueFunc(func(ctx context.Context, newValue string) error {
		if err := c.Set(cacheKey, newValue, s.expiration).Err(); err != nil {
			logrus.Error(err)
		}
		s.mu.Unlock()
		newFnCalled = true
		return foundFn(ctx, newValue, false)
	}))

	if !newFnCalled {
		s.mu.Unlock()
	}

	return err
}

func (s *RedisFetcher) Clear(ctx context.Context, key string) error {
	s.mu.Lock()
	err := s.redis.Del(s.cacheKey(key)).Err()
	s.mu.Unlock()
	return err
}

func (s *RedisFetcher) cacheKey(key string) string {
	return strings.Replace(redisCacheKey, "{key}", key, 1)
}
