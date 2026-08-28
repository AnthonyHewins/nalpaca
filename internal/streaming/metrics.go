package streaming

import (
	"github.com/AnthonyHewins/nalpaca/internal/system"
	"github.com/prometheus/client_golang/prometheus"
)

func newCounter(component, name, help string) prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: system.Name,
		Subsystem: component + "_stream",
		Name:      name,
		Help:      help,
	})
}

type metrics struct {
	publishCount, receiveCount, transformErr, totalErr, marshalErr, pubErr prometheus.Counter
}

func newMetrics(component string) metrics {
	return metrics{
		receiveCount: newCounter(component, "receive_count", "how many messages are received from alpaca for this component"),
		totalErr:     newCounter(component, "total_err", "total error count"),
		marshalErr:   newCounter(component, "marshal_err", "marshal error count"),
		transformErr: newCounter(component, "transform_err", "error count transforming alpaca data types to serializable data"),
		pubErr:       newCounter(component, "pub_err", "nats publish error count"),
		publishCount: newCounter(component, "publish_count", "How many publishes for this component, labeled by symbol"),
	}
}

func (m *metrics) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		m.publishCount,
		m.totalErr,
		m.receiveCount,
		m.marshalErr,
		m.pubErr,
	}
}
