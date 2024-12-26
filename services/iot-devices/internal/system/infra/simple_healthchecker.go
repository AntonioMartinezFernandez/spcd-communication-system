package system_infra

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type SimpleHealthChecker struct {
	tracer trace.Tracer
}

func NewSimpleHealthChecker(tracer trace.Tracer) *SimpleHealthChecker {
	return &SimpleHealthChecker{
		tracer: tracer,
	}
}

func (shc SimpleHealthChecker) Check(ctx context.Context) (map[string]string, error) {
	_, span := shc.tracer.Start(ctx, "HealthChecker.Check")
	defer span.End()

	return map[string]string{
		"system": "ok",
	}, nil
}
