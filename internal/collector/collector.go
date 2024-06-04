package collector

import (
	"log"
	"strings"
)

const (
	DefaultScrapeInterval  = 15
	DefaultMetricsFilePath = "/tmp/node_exporter/textfile"
)

type Collector struct {
	cfg Config
}

type Config struct {
	ScrapeInterval   int
	MetricsFilePath  string
	ConnectFromAllNs bool
	ExcludeSubnets   []string
}

// NewCollector creates a new Collector instance with the given configuration.
func NewCollector(c *Config) *Collector {
	return &Collector{
		cfg: *c,
	}
}

// NewConfig creates a new Config instance with the given fields.
func NewConfig(interval int, output string, exclude string, all bool) *Config {
	excludeSubnets := strings.Split(exclude, ",")

	if len(excludeSubnets) == 1 && excludeSubnets[0] == "" {
		excludeSubnets[0] = "none"
	}
	return &Config{
		ScrapeInterval:   interval,
		MetricsFilePath:  output,
		ConnectFromAllNs: all,
		ExcludeSubnets:   excludeSubnets,
	}
}

// Run starts the metric collection process for the Collector.
func (c *Collector) Run() {
	namespaceInfo := "restrict scraping to connections within the host's network namespace"
	if c.cfg.ConnectFromAllNs {
		namespaceInfo = "gather connections across all network namespaces"
	}

	optionsFormat := "Options:\n" +
		"* metric generation interval: %d\n" +
		"* metric file path: %s\n" +
		"* %s\n" +
		"* exclude patterns for IP addresses: %v\n"
	log.Printf(optionsFormat,
		c.cfg.ScrapeInterval,
		c.cfg.MetricsFilePath,
		namespaceInfo,
		strings.Join(c.cfg.ExcludeSubnets, ", "))

	// for {
	connections, err := c.getConnections()
	if err != nil {
		log.Printf("Error getting connections: %v", err)
		// time.Sleep(time.Duration(c.cfg.ScrapeInterval) * time.Second)
		// continue
	}
	if err := c.writeToFile(c.cfg.MetricsFilePath, connections); err != nil {
		log.Printf("Error writing to file: %v", err)
		// }
		// time.Sleep(time.Duration(c.cfg.ScrapeInterval) * time.Second)
	}
}
