package conf

import (
	"context"
	"runtime/debug"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

//go:generate enumer -type TraceExporter -text -transform lower -trimprefix TraceExporter
type TraceExporter byte

const (
	TraceExporterNone TraceExporter = iota
	TraceExporterStdout
	TraceExporterOTLP
)

type Tracer struct {
	Exporter    TraceExporter `env:"TRACE_EXPORTER" envDefault:"none"`
	ExporterURL string        `env:"TRACE_EXPORTER_URL"`
	Timeout     time.Duration `env:"TRACE_EXPORTER_TIMEOUT" envDefault:"5s"`
}

func (a *Bootstrapper) Tracer(appName string, t *Tracer) (*sdk.TracerProvider, error) {
	var spanExporter sdk.SpanExporter
	var err error

	l := a.Logger.With("config", t)
	switch t.Exporter {
	case TraceExporterStdout:
		spanExporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	case TraceExporterOTLP:
		spanExporter, err = a.otlp(t)
	default:
		l.Info("no tracer specified, creating no-op tracer provider")
		return sdk.NewTracerProvider(sdk.WithBatcher(tracetest.NewNoopExporter())), nil
	}

	if err != nil {
		l.Error("failed creating tracer", "err", err)
		return nil, err
	}

	a.Logger.Info("created tracer")
	return sdk.NewTracerProvider(
		sdk.WithBatcher(spanExporter),
		sdk.WithResource(versionResource(appName)),
	), nil
}

func (a *Bootstrapper) otlp(t *Tracer) (*otlptrace.Exporter, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithReconnectionPeriod(time.Second),
		otlptracegrpc.WithTimeout(t.Timeout),
	}

	if t.ExporterURL != "" {
		opts = append(opts, otlptracegrpc.WithEndpoint(t.ExporterURL))
	}

	return otlptracegrpc.New(context.Background(), opts...)
}

// versionResource returns a resource describing this application.
func versionResource(appName string) *resource.Resource {
	attrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(appName),
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		attrs = append(attrs, semconv.ServiceVersionKey.String(info.Main.Version))
	}

	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL, attrs...),
	)

	return r
}
