package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"goapi/internal/catalog"
	"goapi/internal/report"
	"goapi/pkg/tracing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

func main() {

	cleanup := tracing.InitTracer(serviceName, collectorURL, insecure)
	defer cleanup(context.Background())

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Debug = true
	e.Logger.SetOutput(os.Stdout)

	db, err := gorm.Open(postgres.Open("<actual_dsn>"), &gorm.Config{})
	if err != nil {
		return err
	}

	catalogservice := catalog.NewService(db)
	catalogEndpoint := catalog.NewEndpoint(catalogservice)
	catalogTransport := catalog.NewTransport(catalogEndpoint)
	catalogTransport.RegisterRoutes(e.Group("/api/v1/catalog"))

	reportservice := report.NewService(catalogEndpoint)
	reportEndpoint := report.NewEndpoint(reportservice)
	reportTransport := report.NewTransport(reportEndpoint)
	reportTransport.RegisterRoutes(e.Group("/api/v1/report"))

	e.Use(otelecho.Middleware("goapi-server"))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: os.Stdout,
	}))
	e.Use(middleware.Recover())

	return e.Start(fmt.Sprintf(":%d", 5001))
}
