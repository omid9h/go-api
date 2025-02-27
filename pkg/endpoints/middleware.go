package endpoints

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type ContextKey string

const OtelSpan ContextKey = "OtelSpan"

func RequestCounter[Request, Response any](counter metric.Int64Counter, method string) Middleware[Request, Response] {
	return func(next Endpoint[Request, Response]) Endpoint[Request, Response] {
		return func(ctx context.Context, request Request) (response Response, err error) {
			counter.Add(ctx, 1, metric.WithAttributes(
				attribute.KeyValue{Key: "method", Value: attribute.StringValue(method)},
			))
			return next(ctx, request)
		}
	}
}

func RequestDuration[Request, Response any](histogram metric.Int64Histogram, method string) Middleware[Request, Response] {
	return func(next Endpoint[Request, Response]) Endpoint[Request, Response] {
		return func(ctx context.Context, request Request) (response Response, err error) {
			start := time.Now()
			defer func() {
				histogram.Record(ctx, time.Since(start).Microseconds(), metric.WithAttributes(
					attribute.KeyValue{Key: "method", Value: attribute.StringValue(method)},
				))
			}()
			return next(ctx, request)
		}
	}
}

func OtelTracing[Request, Response any](tracerName, spanName string) Middleware[Request, Response] {
	return func(next Endpoint[Request, Response]) Endpoint[Request, Response] {
		return func(ctx context.Context, request Request) (response Response, err error) {
			tracer := otel.Tracer(tracerName)
			ctx, span := tracer.Start(ctx, spanName)
			defer span.End()
			ctx = context.WithValue(ctx, OtelSpan, span)
			return next(ctx, request)
		}
	}
}
