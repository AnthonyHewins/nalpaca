package conf

import (
	"context"
	"fmt"
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
	TraceExporterStdout TraceExporter = iota
	TraceExporterOTLP
)

type Tracer struct {
	DisableTracing bool          `env:"DISABLE_TRACING" envDefault:"true"`
	Exporter       TraceExporter `env:"TRACE_EXPORTER"`
	ExporterURL    string        `env:"TRACE_EXPORTER_URL"`
	Timeout        time.Duration `env:"TRACE_EXPORTER_TIMEOUT" envDefault:"5s"`
}

func (a *Bootstrapper) Tracer(appName string, t *Tracer) (*sdk.TracerProvider, error) {
	l := a.Logger.With("config", t)

	if t.DisableTracing {
		l.Info("tracing set to disabled, creating noop tracer")
		return sdk.NewTracerProvider(sdk.WithBatcher(tracetest.NewNoopExporter())), nil
	}

	var spanExporter sdk.SpanExporter
	var err error

	switch t.Exporter {
	case TraceExporterStdout:
		spanExporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	case TraceExporterOTLP:
		spanExporter, err = a.otlp(t)
	default:
		l.Error("invalid trace exporter", "val", t.Exporter)
		return nil, fmt.Errorf("invalid trace exporter: %v", t.Exporter)
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
