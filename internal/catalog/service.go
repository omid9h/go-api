package catalog

import (
	"context"
	"log/slog"

	"goapi/pkg/logging"
)

type Service interface {
	ListCatalogs(context.Context, ListCatalogsParams) (ListCatalogsResult, error)
}

type ListCatalogsParams struct {
	Tag string
}

type ListCatalogsResult struct {
	Catalogs []Catalog
}

type Catalog struct {
	ID   string
	Name string
	Tags string
}

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) ListCatalogs(ctx context.Context, params ListCatalogsParams) (result ListCatalogsResult, err error) {

	logger := logging.FromContext(ctx)
	logger.Info("Handling ListCatalogs request", slog.String("tag", params.Tag))

	return
}
