package catalog

import (
	"context"
	"goapi/pkg/endpoints"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var meter = otel.Meter("goapi/internal/catalog")

var apiCounter, _ = meter.Int64Counter(
	"api.counter",
	metric.WithDescription("Number of API calls."),
	metric.WithUnit("{call}"),
)

var apiDuration, _ = meter.Int64Histogram(
	"api.duration",
	metric.WithDescription("Duration of API calls."),
	metric.WithUnit("{microseconds}"),
)

type endpoint struct {
	s Service
}

func (e endpoint) ListCatalogs(ctx context.Context, params ListCatalogsParams) (ListCatalogsResult, error) {
	return endpoints.Chain(
		endpoints.RequestCounter[ListCatalogsParams, ListCatalogsResult](apiCounter, "ListCatalogs"),
		endpoints.RequestDuration[ListCatalogsParams, ListCatalogsResult](apiDuration, "ListCatalogs"),
		endpoints.OtelTracing[ListCatalogsParams, ListCatalogsResult]("goapi/internal/catalog/endpoint", "ListCatalogs"),
	)(e.s.ListCatalogs)(ctx, params)
}

func NewEndpoint(s Service) Service {
	return &endpoint{s}
}
