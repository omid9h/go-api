package catalog

import (
	"context"
	"log/slog"

	"goapi/pkg/logging"

	"gorm.io/gorm"
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

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *service {
	return &service{db: db}
}

func (s *service) ListCatalogs(ctx context.Context, params ListCatalogsParams) (result ListCatalogsResult, err error) {
	db := s.db.WithContext(ctx) // context aware gorm session

	_ = db // use db to perform database operations in real example

	logger := logging.FromContext(ctx)
	logger.Info("Handling ListCatalogs request", slog.String("tag", params.Tag))

	return
}
