package strcache

import (
	"context"
)

type FoundFunc func(ctx context.Context, value string, isNew bool) error
type NotFoundFunc func(ctx context.Context, newFn NewValueFunc) error
type NewValueFunc func(ctx context.Context, newValue string) error

type Fetcher interface {
	Fetch(ctx context.Context, key string, notFoundFn NotFoundFunc, foundFn FoundFunc) error
	Clear(ctx context.Context, key string) error
}
