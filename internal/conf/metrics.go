package conf

import (
	"github.com/AnthonyHewins/nalpaca/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Creates a prometheus metric HTTP server. Pass a non-nil logger to log errors. By default this will automatically
// create a version gaugevec collector and append it to your server. Pass any other prom collectors into this function
// to track other metrics
func (b *Bootstrapper) PrometheusHTTP(m *metrics.PromConfig, collectors ...prometheus.Collector) (metrics.Prom, error) {
	if m.DisableMetrics {
		b.Logger.Info("metrics disabled, not creating prom metrics")
		return metrics.Prom{}, nil
	}

	return metrics.NewProm(b.Logger, m)
}
