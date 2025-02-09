package catalog

import (
	"context"
	"goapi/pkg/endpoint"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var meter = otel.Meter("goapi/internal/catalog")

type Endpoint struct {
	ListCatalogs func(context.Context, ListCatalogsParams) (ListCatalogsResult, error)
}

func NewEndpoint(s Service) (endpoints Endpoint) {
	endpoints = Endpoint{
		ListCatalogs: endpoint.MakeEndpoint(s.ListCatalogs),
	}

	apiCounter, _ := meter.Int64Counter(
		"api.counter",
		metric.WithDescription("Number of API calls."),
		metric.WithUnit("{call}"),
	)

	apiDuration, _ := meter.Int64Histogram(
		"api.duration",
		metric.WithDescription("Duration of API calls."),
		metric.WithUnit("{microseconds}"),
	)

	endpoints.ListCatalogs = endpoint.Chain(
		endpoint.RequestCounter[ListCatalogsParams, ListCatalogsResult](apiCounter, "ListCatalogs"),
		endpoint.RequestDuration[ListCatalogsParams, ListCatalogsResult](apiDuration, "ListCatalogs"),
	)(endpoints.ListCatalogs)

	return
}
