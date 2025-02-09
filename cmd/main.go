package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"goapi/internal/catalog"
	"goapi/pkg/tracing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
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

	catalogservice := catalog.NewService()
	catalogEndpoint := catalog.NewEndpoint(catalogservice)
	catalogTransport := catalog.NewTransport(catalogEndpoint)
	catalogTransport.RegisterRoutes(e.Group("/api/v1/catalog"))

	e.Use(otelecho.Middleware("goapi-server"))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: os.Stdout,
	}))
	e.Use(middleware.Recover())

	e.Start(fmt.Sprintf(":%d", 5001))
	return nil
}
