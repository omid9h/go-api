package logging

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

// baseLogger is our global logger, set to output to Stdout in text format.
// Adjust as needed (e.g., JSON, severity levels, etc.)
var baseLogger *slog.Logger

func init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	baseLogger = slog.New(handler)
}

// FromContext returns a logger that includes OpenTelemetry trace IDs in its fields,
// enabling log and trace correlation in SigNoz or any OT-compatible backend.
func FromContext(ctx context.Context) *slog.Logger {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return baseLogger
	}
	return baseLogger.With(
		slog.String("trace_id", spanCtx.TraceID().String()),
		slog.String("span_id", spanCtx.SpanID().String()),
	)
}
