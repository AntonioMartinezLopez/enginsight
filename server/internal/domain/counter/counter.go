package counter

import (
	"context"
	"unicode/utf8"
)

type service struct {
	store CounterStore
}

func New(store CounterStore) Counter {
	return &service{store: store}
}

func (c *service) Count(ctx context.Context, message string) (int, error) {
	if err := c.store.Increment(ctx); err != nil {
		return 0, NewInternalError("failed to increment message count", err)
	}
	return utf8.RuneCountInString(message), nil
}

func (c *service) GetNumberOfProcessedMessages(ctx context.Context) (Count, error) {
	count, err := c.store.GetMessageCount(ctx)
	if err != nil {
		return 0, NewInternalError("failed to get message count", err)
	}
	return count, nil
}
