package counter

import "context"

type Counter interface {
	Count(ctx context.Context, message string) (int, error)
	GetNumberOfProcessedMessages(ctx context.Context) (Count, error)
}

type CounterStore interface {
	Increment(ctx context.Context) error
	GetMessageCount(ctx context.Context) (Count, error)
}
