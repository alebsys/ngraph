package exporter

import (
	"log"

	"github.com/alebsys/ngraph/internal/collector"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	labels        = []string{"src_ip", "dst_ip"}
	incomingConns = prometheus.NewDesc(
		"network_connections_incoming_total",
		"Total number of incoming network connections between source and destination IP addresses.",
		labels, nil,
	)
	outgoingConns = prometheus.NewDesc(
		"network_connections_outgoing_total",
		"Total number of outgoing network connections between source and destination IP addresses.",
		labels, nil,
	)
)

type Exporter struct {
	cfg       Config
	collector *collector.Collector
}

type Config struct {
	ListeningAddress string
	MetricsEndpoint  string
}

// NewCollector creates a new Exporter instance with the given configuration.
func NewExporter(c Config, coll *collector.Collector) *Exporter {
	return &Exporter{
		cfg:       c,
		collector: coll,
	}
}

// NewConfig creates a new Config instance for Exporter with the given fields.
func NewConfig(addr, endpoint string) *Config {
	return &Config{
		ListeningAddress: addr,
		MetricsEndpoint:  endpoint,
	}
}

// TODO: ...implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	conns, err := e.collector.GetConnections()
	if err != nil {
		log.Printf("Error getting connections: %v", err)
	} else {
		for conn, v := range conns {
			if conn.Direction == "output" {
				ch <- prometheus.MustNewConstMetric(outgoingConns, prometheus.GaugeValue, v, conn.SrcIP, conn.DstIP)
			} else {
				ch <- prometheus.MustNewConstMetric(incomingConns, prometheus.GaugeValue, v, conn.DstIP, conn.SrcIP)
			}
		}
	}
}

// TODO: ...implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- incomingConns
	ch <- outgoingConns
}
