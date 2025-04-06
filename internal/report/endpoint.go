package report

import (
	"context"
	"goapi/internal/catalog"
	"goapi/pkg/endpoints"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var meter = otel.Meter("goapi/internal/report")

var apiCounter, _ = meter.Int64Counter(
	"api.counter",
	metric.WithDescription("Number of API calls."),
	metric.WithUnit("{call}"),
)

type endpoint struct {
	s Service
}

func (e endpoint) ReportCatalogs(ctx context.Context, params catalog.ListCatalogsParams) (catalog.ListCatalogsResult, error) {
	return endpoints.Chain(
		endpoints.RequestCounter[catalog.ListCatalogsParams, catalog.ListCatalogsResult](apiCounter, "ReportCatalogs"),
	)(e.s.ReportCatalogs)(ctx, params)
}

func NewEndpoint(s Service) Service {
	return &endpoint{s}
}
