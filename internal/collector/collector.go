package collector

import (
	"log"
	"time"
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
}

// NewCollector creates a new Collector instance with the given configuration.
func NewCollector(c *Config) *Collector {
	return &Collector{
		cfg: *c,
	}
}

// Run starts the metric collection process for the Collector.
func (c *Collector) Run() {
	namespaceInfo := "only from host network namespace"
	if c.cfg.ConnectFromAllNs {
		namespaceInfo = "from all network namespaces"
	}
	log.Printf("Options:\n* interval for generating metrics: %d\n* path for metric files: %s\n* scrape connections %s\n", c.cfg.ScrapeInterval, c.cfg.MetricsFilePath, namespaceInfo)

	for {
		connections, err := c.getConnections()
		if err != nil {
			log.Println(err)
			continue
		}
		if err := c.writeToFile(c.cfg.MetricsFilePath, connections); err != nil {
			log.Println(err)
			continue
		}
		time.Sleep(time.Duration(c.cfg.ScrapeInterval) * time.Second)
	}
}
