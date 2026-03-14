package jrpc

import "context"

type CountService interface {
	Count(ctx context.Context, message string) (int, error)
}
