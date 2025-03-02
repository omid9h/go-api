package report

import (
	"context"
	"goapi/internal/catalog"
	"log/slog"

	"goapi/pkg/logging"
)

type Service interface {
	ReportCatalogs(context.Context, catalog.ListCatalogsParams) (catalog.ListCatalogsResult, error)
}

type service struct {
	catalogService catalog.Service
}

func NewService(catalogService catalog.Service) *service {
	return &service{catalogService: catalogService}
}

func (s *service) ReportCatalogs(ctx context.Context, params catalog.ListCatalogsParams) (result catalog.ListCatalogsResult, err error) {
	logger := logging.FromContext(ctx)
	logger.Info("Handling ReportCatalogs request", slog.String("tag", params.Tag))

	result, err = s.catalogService.ListCatalogs(ctx, params)
	if err != nil {
		return result, err
	}

	return result, nil
}
