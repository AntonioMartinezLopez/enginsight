package jrpc

import (
	"context"
	"net/http"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/jhttp"
)

type Client interface {
	CountService
	Close() error
}

type client struct {
	cli *jrpc2.Client
}

type ClientOptions struct {
	// Custom HTTP client (optional), e.g. for custom timeouts or tracing
	HTTPClient *http.Client
}

func NewClient(url string, opts ClientOptions) (Client, error) {
	channelOpts := &jhttp.ChannelOptions{}

	if opts.HTTPClient != nil {
		channelOpts.Client = opts.HTTPClient
	}

	channel := jhttp.NewChannel(url, channelOpts)
	cli := jrpc2.NewClient(channel, nil)

	return &client{cli: cli}, nil
}

func (c *client) Count(ctx context.Context, message string) (int, error) {
	req := CountRequest{Message: message}
	var response CountResponse
	err := c.cli.CallResult(ctx, MethodCount, req, &response)
	return response.Count, err
}

func (c *client) Close() error {
	return c.cli.Close()
}
