package jrpc

import (
	"context"
	"net/http"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/jhttp"
)

type Server struct {
	bridge jhttp.Bridge
}

func (s *Server) Close() error {
	return s.bridge.Close()
}

// ServeHTTP implements http.Handler, allowing Server to be used directly with HTTP servers.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.bridge.ServeHTTP(w, r)
}

type Handlers struct {
	CountService CountService
}

func NewServer(handlers Handlers) *Server {
	mux := handler.Map{}

	if handlers.CountService != nil {
		mux[MethodCount] = handler.New(func(ctx context.Context, req CountRequest) (CountResponse, error) {
			count, err := handlers.CountService.Count(ctx, req.Message)
			if err != nil {
				return CountResponse{}, err
			}
			return CountResponse{Count: count}, nil
		})
	}

	serverOpts := &jrpc2.ServerOptions{}
	serverOpts.Logger = jrpc2.StdLogger(nil)
	bridge := jhttp.NewBridge(mux, &jhttp.BridgeOptions{
		Server: serverOpts,
	})

	return &Server{
		bridge: bridge,
	}
}
