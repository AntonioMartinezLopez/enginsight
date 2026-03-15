package counter_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"unicode/utf8"

	pkgerrors "github.com/AntonioMartinezLopez/enginsight/pkg"
	"github.com/AntonioMartinezLopez/enginsight/server/internal/domain/counter"
	"github.com/AntonioMartinezLopez/enginsight/server/internal/domain/counter/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCounter(t *testing.T) {
	if !testing.Short() {
		t.Skip()
	}

	tests := []struct {
		name     string
		message  string
		expected int
	}{
		{
			name:     "empty message with 0 characters",
			message:  "",
			expected: 0,
		},
		{
			name:     "small message with few characters",
			message:  "Hi",
			expected: 2,
		},
		{
			name:     "message with newline escape character",
			message:  `Hello\nWorld`,
			expected: 12,
		},
		{
			name:     "message with tab escape character",
			message:  `Hello\tWorld`,
			expected: 12,
		},
		{
			name:     "message with various escape characters",
			message:  `Line1\nLine2\tTabbed\rCarriage`,
			expected: 30,
		},
		{
			name:     "message with special unicode characters",
			message:  "Hello 🌍 World 🚀",
			expected: 15,
		},
		{
			name:     "message with quotes and backslashes",
			message:  `Quote: "test" and backslash: \path\to\file`,
			expected: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock CounterStore
			mockStore := mocks.NewMockCounterStore(t)

			// Set up expectations for the mock store
			mockStore.EXPECT().Increment(mock.Anything).Return(nil).Times(1)

			// Create the Counter service with the mock store
			counterService := counter.New(mockStore)

			// Call the Count method
			count, err := counterService.Count(context.Background(), tt.message)
			require.NoError(t, err)
			require.Equal(t, tt.expected, count)
		})
	}
}

func TestCounterProcessMultipleMessages(t *testing.T) {
	if !testing.Short() {
		t.Skip()
	}

	// Create a mock CounterStore
	mockStore := mocks.NewMockCounterStore(t)

	// Set up expectations for the mock store
	mockStore.EXPECT().Increment(mock.Anything).Return(nil).Times(3)
	mockStore.EXPECT().GetMessageCount(mock.Anything).Return(counter.Count(3), nil).Times(1)

	// Create the Counter service with the mock store
	counterService := counter.New(mockStore)

	// Process multiple messages concurrently to test thread safety
	messages := []string{"Hello", "World", "!"}
	wg := sync.WaitGroup{}
	wg.Add(len(messages))

	for _, msg := range messages {
		go func(m string) {
			count, err := counterService.Count(context.Background(), m)
			require.NoError(t, err)
			require.Equal(t, utf8.RuneCountInString(m), count)
			wg.Done()
		}(msg)
	}
	wg.Wait()

	// Call the GetNumberOfProcessedMessages method
	totalCount, err := counterService.GetNumberOfProcessedMessages(context.Background())
	require.NoError(t, err)
	require.Equal(t, counter.Count(3), totalCount)
}

func TestCounterIncrementError(t *testing.T) {
	if !testing.Short() {
		t.Skip()
	}

	// Create a mock CounterStore
	mockStore := mocks.NewMockCounterStore(t)

	// Set up expectations - Increment returns a generic error
	storeErr := errors.New("store failure")
	mockStore.EXPECT().Increment(mock.Anything).Return(storeErr).Times(1)

	// Create the Counter service with the mock store
	counterService := counter.New(mockStore)

	// Call the Count method - should return an error
	count, err := counterService.Count(context.Background(), "test message")
	require.Error(t, err)
	require.Equal(t, 0, count)

	// Assert the error is of type InternalError from counter package
	var internalErr *pkgerrors.Error
	require.ErrorAs(t, err, &internalErr)
	require.Equal(t, counter.ErrCodeInternal, internalErr.Code)
	require.Contains(t, internalErr.Error(), "failed to increment message count")
	require.ErrorIs(t, internalErr.Unwrap(), storeErr)
}

func TestCounterGetMessageCountError(t *testing.T) {
	if !testing.Short() {
		t.Skip()
	}

	// Create a mock CounterStore
	mockStore := mocks.NewMockCounterStore(t)

	// Set up expectations - GetMessageCount returns a generic error
	storeErr := errors.New("store failure")
	mockStore.EXPECT().GetMessageCount(mock.Anything).Return(counter.Count(0), storeErr).Times(1)

	// Create the Counter service with the mock store
	counterService := counter.New(mockStore)

	// Call the GetNumberOfProcessedMessages method - should return an error
	count, err := counterService.GetNumberOfProcessedMessages(context.Background())
	require.Error(t, err)
	require.Equal(t, counter.Count(0), count)

	// Assert the error is of type InternalError from counter package
	var internalErr *pkgerrors.Error
	require.ErrorAs(t, err, &internalErr)
	require.Equal(t, counter.ErrCodeInternal, internalErr.Code)
	require.Contains(t, internalErr.Error(), "failed to get message count")
	require.ErrorIs(t, internalErr.Unwrap(), storeErr)
}
