package test

import (
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/AntonioMartinezLopez/enginsight/jrpc"
	"github.com/AntonioMartinezLopez/enginsight/server/internal/domain/counter"
	"github.com/AntonioMartinezLopez/enginsight/server/internal/infra/rpc"
	"github.com/AntonioMartinezLopez/enginsight/server/internal/infra/store"
)

// testServerSetup contains the test server components
type testServerSetup struct {
	counterService counter.Counter
	httpServer     *httptest.Server
	rpcServer      *jrpc.Server
}

// setupTestServer creates and configures a test server with all dependencies
func setupTestServer(t *testing.T) *testServerSetup {
	t.Helper()

	// Set up application
	store := store.New()
	counterService := counter.New(store)
	rpc := rpc.New(rpc.ServerConfig{
		Service: counterService,
	})

	// Set up HTTP test server with the JSON-RPC handler
	httpServer := httptest.NewServer(rpc.Handler())

	return &testServerSetup{
		counterService: counterService,
		httpServer:     httpServer,
		rpcServer:      rpc.Handler(),
	}
}

// Close cleans up all server resources
func (s *testServerSetup) Close() {
	if s.httpServer != nil {
		s.httpServer.Close()
	}
	if s.rpcServer != nil {
		s.rpcServer.Close()
	}
}

func TestCounterIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// Set up test server
	setup := setupTestServer(t)
	defer setup.Close()

	// Set up client to connect to the test server
	client, err := jrpc.NewClient(setup.httpServer.URL, jrpc.ClientOptions{})
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

func TestCounterIntegrationMultipleClients(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// Set up test server
	setup := setupTestServer(t)
	defer setup.Close()

	// Set up multiple clients to connect to the test server
	client1, err := jrpc.NewClient(setup.httpServer.URL, jrpc.ClientOptions{})
	if err != nil {
		t.Fatalf("Failed to create client1: %v", err)
	}
	defer client1.Close()

	client2, err := jrpc.NewClient(setup.httpServer.URL, jrpc.ClientOptions{})
	if err != nil {
		t.Fatalf("Failed to create client2: %v", err)
	}
	defer client2.Close()

	// Test cases
	tests := []struct {
		name           string
		message        string
		expectedLength int
	}{
		{
			name:           "empty message",
			message:        "",
			expectedLength: 0,
		},
		{
			name:           "simple message",
			message:        "Hello",
			expectedLength: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := sync.WaitGroup{}
			wg.Add(2)

			// Launch concurrent requests from both clients
			go func() {
				defer wg.Done()
				result, err := client1.Count(t.Context(), tt.message)
				if err != nil {
					t.Errorf("Client1 Count failed: %v", err)
					return
				}
				if result != tt.expectedLength {
					t.Errorf("Client1 expected length %d, got %d", tt.expectedLength, result)
				}
			}()

			go func() {
				defer wg.Done()
				result, err := client2.Count(t.Context(), tt.message)
				if err != nil {
					t.Errorf("Client2 Count failed: %v", err)
					return
				}
				if result != tt.expectedLength {
					t.Errorf("Client2 expected length %d, got %d", tt.expectedLength, result)
				}
			}()

			// Wait for both goroutines to complete before subtest exits
			wg.Wait()
		})
	}

	// Final verification of total count
	totalCount, _ := setup.counterService.GetNumberOfProcessedMessages(t.Context())
	expectedTotal := counter.Count(len(tests) * 2) // Each test case is processed by 2 clients
	if totalCount != expectedTotal {
		t.Errorf("Expected total count %d, got %d", expectedTotal, totalCount)
	}
}
