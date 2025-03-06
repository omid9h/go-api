package endpoints

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type ContextKey string

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
			var span trace.Span

			// Check if there's an existing span in the context
			parentSpan := trace.SpanFromContext(ctx)
			if parentSpan.SpanContext().IsValid() {
				// Create a new span as a child of the existing span
				ctx, span = tracer.Start(ctx, spanName, trace.WithLinks(trace.LinkFromContext(ctx)))
			} else {
				// No existing span, start a new root span
				ctx, span = tracer.Start(ctx, spanName)
			}

			defer span.End()

			return next(ctx, request)
		}
	}
}
