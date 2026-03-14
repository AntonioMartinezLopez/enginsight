package counter

import "context"

type service struct {
	store CounterStore
}

func New(store CounterStore) Counter {
	return &service{store: store}
}

func (c *service) Count(ctx context.Context, message string) (int, error) {
	if err := c.store.Increment(ctx); err != nil {
		return 0, err
	}
	return len(message), nil
}

func (c *service) GetNumberOfProcessedMessages(ctx context.Context) (Count, error) {
	count, err := c.store.GetMessageCount(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}
