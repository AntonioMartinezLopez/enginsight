package test

import (
	"net/http/httptest"
	"testing"

	"github.com/AntonioMartinezLopez/enginsight/jrpc"
	"github.com/AntonioMartinezLopez/enginsight/server/internal/domain/counter"
	"github.com/AntonioMartinezLopez/enginsight/server/internal/infra/store"
)

func TestCounterIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// set up application
	store := store.New()
	counterService := counter.New(store)

	// Create the JSON-RPC server with our service implementation
	handlers := jrpc.Handlers{
		CountService: counterService,
	}
	rpcServer := jrpc.NewServer(handlers)
	defer rpcServer.Close()

	// Set up HTTP test server with the JSON-RPC handler
	httpServer := httptest.NewServer(rpcServer)
	defer httpServer.Close()

	// Set up client to connect to the test server
	client, err := jrpc.NewClient(httpServer.URL, jrpc.ClientOptions{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Test cases
	tests := []struct {
		name     string
		message  string
		expected int
	}{
		{
			name:     "empty message",
			message:  "",
			expected: 0,
		},
		{
			name:     "simple message",
			message:  "Hello",
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := client.Count(t.Context(), tt.message)
			if err != nil {
				t.Fatalf("Count failed: %v", err)
			}
			if count != tt.expected {
				t.Errorf("Expected count %d, got %d", tt.expected, count)
			}
		})
	}

}
