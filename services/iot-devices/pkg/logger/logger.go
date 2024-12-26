package logger

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

type Logger interface {
	Error(ctx context.Context, message string, items ...slog.Attr)
	Warn(ctx context.Context, message string, items ...slog.Attr)
	Info(ctx context.Context, message string, items ...slog.Attr)
	Debug(ctx context.Context, message string, items ...slog.Attr)
}

type logger struct {
	sLogger *slog.Logger
}

func NewLogger(level string) *logger {
	opts := &slog.HandlerOptions{Level: levelFromStringLevel(level)}

	return &logger{
		sLogger: slog.New(NewPrettyLogHandler(opts).WithGroup("data")),
	}
}

func NewJsonLogger(level string) *logger {
	jsonHandler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: levelFromStringLevel(level),
		},
	)

	return &logger{
		sLogger: slog.New(jsonHandler),
	}
}

func NewOtelInstrumentalizedLogger(level string) *logger {
	jsonHandler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: levelFromStringLevel(level),
		},
	)

	instrumentalizedHandler := handlerWithSpanContext(jsonHandler)

	return &logger{
		sLogger: slog.New(instrumentalizedHandler),
	}
}

func (l *logger) Error(ctx context.Context, message string, items ...slog.Attr) {
	l.sLogger.LogAttrs(ctx, slog.LevelError, message, items...)
}

func (l *logger) Warn(ctx context.Context, message string, items ...slog.Attr) {
	l.sLogger.LogAttrs(ctx, slog.LevelWarn, message, items...)
}

func (l *logger) Info(ctx context.Context, message string, items ...slog.Attr) {
	l.sLogger.LogAttrs(ctx, slog.LevelInfo, message, items...)
}

func (l *logger) Debug(ctx context.Context, message string, items ...slog.Attr) {
	l.sLogger.LogAttrs(ctx, slog.LevelDebug, message, items...)
}

func levelFromStringLevel(lvl string) slog.Level {
	var logLevel slog.Level
	switch level := lvl; level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelWarn
	}

	return logLevel
}

func handlerWithSpanContext(handler slog.Handler) *spanContextLogHandler {
	return &spanContextLogHandler{Handler: handler}
}

type spanContextLogHandler struct {
	slog.Handler
}

func (slh *spanContextLogHandler) Handle(ctx context.Context, record slog.Record) error {
	// Get the SpanContext from the golang Context.
	if spanCtx := trace.SpanContextFromContext(ctx); spanCtx.IsValid() {
		// Add trace context attributes following Cloud Logging structured log format described
		// in https://cloud.google.com/logging/docs/structured-logging#special-payload-fields
		record.AddAttrs(
			slog.Any("trace", spanCtx.TraceID()),
		)
		record.AddAttrs(
			slog.Any("spanId", spanCtx.SpanID()),
		)
		record.AddAttrs(
			slog.Bool("traceSampled", spanCtx.TraceFlags().IsSampled()),
		)
	}
	return slh.Handler.Handle(ctx, record)
}
