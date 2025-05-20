package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:generate enumer -type MetricsErrPolicy -text -trimprefix MetricsErrPolicy -transform lower
type MetricsErrPolicy byte

const (
	// If errs are encountereed in metrics, return them in HTTP to the listener
	MetricsErrPolicyReturn MetricsErrPolicy = iota

	// If errs are encountered in metrics, continue past them
	MetricsErrPolicyContinue

	// If errs are encountered in metrics, panic
	MetricsErrPolicyPanic
)

type Metrics struct {
	// Disable the prometheus metrics server.
	DisableMetrics bool `env:"DISABLE_METRICS" envDefault:"true"`

	// Port to expose for the Prometheus HTTP metrics API.
	HTTPMetricsPort uint16 `env:"METRICS_PORT" envDefault:"8088"`

	HTTPMetricsTimeout time.Duration `env:"METRICS_WRITE_TIMEOUT" envDefault:"10s"`

	// The policy to use when an error is encountered; you can specify:
	// return: send the error back to the requestor
	// continue: move past the error
	// panic: panic on error
	HTTPMetricsErrPolicy MetricsErrPolicy `env:"METRICS_ERR_POLICY" envDefault:""`

	HTTPMetricsMaxRequestsInFlight uint `env:"METRICS_MAX_REQ_IN_FLIGHT" envDefault:""`
}

type promLogger struct {
	l *slog.Logger
}

func (p promLogger) Println(entries ...any) {
	for _, v := range entries {
		switch x := v.(type) {
		case prometheus.MultiError:
			for _, err := range x {
				p.l.Error(err.Error())
			}
		case fmt.Stringer:
			p.l.Error(x.String())
		default:
			p.l.Error(fmt.Sprintf("%+v", x))
		}
	}
}

// Creates a prometheus metric HTTP server. Pass a non-nil logger to log errors. By default this will automatically
// create a version gaugevec collector and append it to your server. Pass any other prom collectors into this function
// to track other metrics
func (b *Bootstrapper) PrometheusHTTP(m *Metrics, collectors ...prometheus.Collector) (*http.Server, error) {
	if m.DisableMetrics {
		b.Logger.Info("metrics disabled, not creating prom metrics")
		return nil, nil
	}

	logger := b.Logger

	reg := prometheus.NewRegistry()

	listenAddr := fmt.Sprintf(":%d", m.HTTPMetricsPort)

	h := http.NewServeMux()
	h.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		ErrorLog:            promLogger{logger},
		ErrorHandling:       promhttp.HandlerErrorHandling(m.HTTPMetricsErrPolicy),
		Registry:            reg,
		MaxRequestsInFlight: int(m.HTTPMetricsMaxRequestsInFlight),
		Timeout:             m.HTTPMetricsTimeout,
	}))

	info, ok := debug.ReadBuildInfo()
	if !ok {
		msg := "failed reading build info: your go binary is not built in module mode"
		logger.Error(msg)
		return nil, errors.New(msg)
	}

	versionGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "apiserver",
		Name:      "version",
		Help:      "App version",
	}, []string{"version"})

	http.HandleFunc("/version", func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.MarshalIndent(info, "", " ")
		w.Write(b)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	for _, collector := range append(collectors, versionGauge) {
		if err := reg.Register(collector); err != nil {
			return nil, err
		}
	}

	logger.Info("created metrics server", "conf", m)
	return &http.Server{
		Addr:              listenAddr,
		ReadTimeout:       m.HTTPMetricsTimeout,
		ReadHeaderTimeout: m.HTTPMetricsTimeout,
		WriteTimeout:      m.HTTPMetricsTimeout,
	}, nil
}
