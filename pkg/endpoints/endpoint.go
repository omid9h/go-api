package endpoints

import (
	"context"
)

type Endpoint[Request, Response any] func(context.Context, Request) (Response, error)

func Nop(context.Context, interface{}) (interface{}, error) { return struct{}{}, nil }

type Middleware[Request, Response any] func(Endpoint[Request, Response]) Endpoint[Request, Response]

func Chain[Request, Response any](outer Middleware[Request, Response], others ...Middleware[Request, Response]) Middleware[Request, Response] {
	return func(next Endpoint[Request, Response]) Endpoint[Request, Response] {
		for i := len(others) - 1; i >= 0; i-- {
			next = others[i](next)
		}
		return outer(next)
	}
}

func MakeEndpoint[Request, Response any](e func(context.Context, Request) (Response, error)) Endpoint[Request, Response] {
	return func(ctx context.Context, request Request) (Response, error) {
		return e(ctx, request)
	}
}

type NoResult struct{}
