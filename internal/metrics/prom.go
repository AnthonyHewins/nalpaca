package metrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/system"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prom struct {
	logger *slog.Logger
	prometheus.Registerer
	Server *http.Server
}

func (p Prom) Register(m ...prometheus.Collector) error {
	if p.Registerer == nil {
		return nil
	}

	for _, v := range m {
		if err := p.Registerer.Register(v); err != nil {
			p.logger.Error("failed registering metric", "err", err)
			return err
		}
	}

	return nil
}

type promLogger struct{ l *slog.Logger }

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

//go:generate enumer -type MetricsErrPolicy -text -trimprefix MetricsErrPolicy
type MetricsErrPolicy byte

const (
	// If errs are encountereed in metrics, return them in HTTP to the listener
	MetricsErrPolicyReturn MetricsErrPolicy = iota

	// If errs are encountered in metrics, continue past them
	MetricsErrPolicyContinue

	// If errs are encountered in metrics, panic
	MetricsErrPolicyPanic
)

type PromConfig struct {
	// Disable the prometheus metrics server.
	DisableMetrics bool `env:"DISABLE_METRICS" envDefault:"true" desc:"Disable the prometheus metrics server."`

	// Port to expose for the Prometheus HTTP metrics API.
	HTTPMetricsPort uint16 `env:"METRICS_PORT" envDefault:"8088" desc:"Port to expose for the Prometheus HTTP metrics API."`

	HTTPMetricsTimeout time.Duration `env:"METRICS_WRITE_TIMEOUT" envDefault:"10s" desc:"Write timeout for the Prometheus HTTP metrics server"`

	// The policy to use when an error is encountered; you can specify:
	// return: send the error back to the requestor
	// continue: move past the error
	// panic: panic on error
	HTTPMetricsErrPolicy MetricsErrPolicy `env:"METRICS_ERR_POLICY" envDefault:"" desc:"How to handle an error gathering metrics: return it to the requestor, continue past it, or panic"`

	HTTPMetricsMaxRequestsInFlight uint `env:"METRICS_MAX_REQ_IN_FLIGHT" envDefault:"" desc:"Maximum concurrent requests the metrics HTTP handler will serve; 0 disables the limit"`
}

func NewProm(logger *slog.Logger, m *PromConfig, initialCollectors ...prometheus.Collector) (Prom, error) {
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
		return Prom{}, errors.New(msg)
	}

	for _, collector := range initialCollectors {
		if err := reg.Register(collector); err != nil {
			logger.Error("failed registering metric", "err", err)
			return Prom{}, err
		}
	}

	for _, v := range [...]struct {
		name, value string
		vec         *prometheus.GaugeVec
	}{
		{
			"version",
			system.Version,
			GaugeVec("system", "version", "App version", []string{"version"}),
		},
		{
			"commit",
			system.Commit,
			GaugeVec("system", "commit", "Commit for this build", []string{"commit"}),
		},
		{
			"build time",
			system.BuildTime,
			GaugeVec("system", "build_time", "When the binary was built", []string{"time"}),
		},
	} {
		v.vec.WithLabelValues(v.value).Set(1)
		if err := reg.Register(v.vec); err != nil {
			logger.Error("failed registering system level metric", "err", err, "name", v.name)
			return Prom{}, err
		}
	}

	h.HandleFunc("/version", func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.MarshalIndent(map[string]any{
			"build":     info,
			"commit":    system.Commit,
			"buildTime": system.BuildTime,
			"version":   system.Version,
		}, "", " ")
		w.Write(b)
	})

	h.HandleFunc("/healthz", func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	logger.Info("created metrics server", "conf", m)
	return Prom{
		Registerer: reg,
		Server: &http.Server{
			Addr:              listenAddr,
			Handler:           h,
			ReadTimeout:       m.HTTPMetricsTimeout,
			ReadHeaderTimeout: m.HTTPMetricsTimeout,
			WriteTimeout:      m.HTTPMetricsTimeout,
		},
	}, nil
}

type PromCollector interface{ Metrics() []prometheus.Collector }

func Counter(subsystem, name, help string) prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: system.Name,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	})
}

func GaugeVec(subsystem, name, help string, labels []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: system.Name,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	}, labels)
}
