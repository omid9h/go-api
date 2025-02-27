package catalog

import (
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
	g.GET("/catalogs", t.ListCatalogs)
}

func (t *transport) ListCatalogs(c echo.Context) error {
	params := ListCatalogsParams{}
	if err := c.Bind(&params); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	result, err := t.endpoints.ListCatalogs(c.Request().Context(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
