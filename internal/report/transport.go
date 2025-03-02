package report

import (
	"goapi/internal/catalog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type transport struct {
	endpoints Service
}

func NewTransport(endpoints Service) *transport {
	return &transport{
		endpoints: endpoints,
	}
}

func (t *transport) RegisterRoutes(g *echo.Group) {
	g.GET("/catalogs", t.ReportCatalogs)
}

func (t *transport) ReportCatalogs(c echo.Context) error {
	params := catalog.ListCatalogsParams{}
	if err := c.Bind(&params); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	result, err := t.endpoints.ReportCatalogs(c.Request().Context(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
