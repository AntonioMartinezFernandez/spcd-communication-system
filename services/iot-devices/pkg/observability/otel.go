package observability

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel"
	otel_metric "go.opentelemetry.io/otel/metric"
	otel_semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	otel_trace "go.opentelemetry.io/otel/trace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const schemaName = "go.opentelemetry.io/otel"

type OtelObservability struct {
	Tracer       otel_trace.Tracer
	Meter        otel_metric.Meter
	ShutdownFunc func(context.Context) error
}

func InitOpenTelemetryObservability(
	ctx context.Context,
	grpcConnection *grpc.ClientConn,
	serviceName string,
	serviceVersion string,
) (*OtelObservability, error) {
	var shutdownFuncs []func(context.Context) error
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	otelObservability := &OtelObservability{
		Tracer:       otel.Tracer(schemaName),
		Meter:        otel.Meter(schemaName),
		ShutdownFunc: shutdown,
	}

	handleErr := func(inErr error) error {
		err := errors.Join(inErr, shutdown(ctx))
		return err
	}

	// Setup propagator
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Setup resource
	resource := resource.NewWithAttributes(
		otel_semconv.SchemaURL,
		otel_semconv.ServiceName(serviceName),
		otel_semconv.ServiceVersion(serviceVersion),
	)

	// Setup trace provider
	tracerProvider, err := newTraceProvider(ctx, grpcConnection, resource)
	if err != nil {
		e := handleErr(err)
		return otelObservability, e
	}

	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Setup meter provider
	meterProvider, err := newMeterProvider(ctx, grpcConnection, resource)
	if err != nil {
		e := handleErr(err)
		return otelObservability, e
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return otelObservability, nil
}

func InitGrpcConnInsecure(host, port string) (*grpc.ClientConn, error) {
	return grpc.NewClient(fmt.Sprintf("%s:%s", host, port),
		// Note: TLS is recommended in production
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(ctx context.Context, grpcConnection *grpc.ClientConn, resource *resource.Resource) (*trace.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithGRPCConn(grpcConnection),
	)
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithResource(resource),
		trace.WithBatcher(traceExporter,
			trace.WithBatchTimeout(3*time.Second), // Default is 5s
		),
	)
	return traceProvider, nil
}

func newMeterProvider(ctx context.Context, grpcConnection *grpc.ClientConn, resource *resource.Resource) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithGRPCConn(grpcConnection),
	)
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(resource),
		metric.WithReader(
			metric.NewPeriodicReader(
				metricExporter,
				metric.WithInterval(30*time.Second), // Default is 1m
			),
		),
	)
	return meterProvider, nil
}
